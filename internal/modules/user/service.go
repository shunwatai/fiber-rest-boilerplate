package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-api-starter/internal/auth"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{r}
}

var cfg = config.Cfg

/* this func for generate the jwt claims like the access & refresh tokens */
func GenerateUserToken(user User, tokenType string) *jwt.Token {
	var expireTime = time.Now().Add(time.Minute * 10).Unix() // 10 mins for access token?

	cfg.LoadEnvVariables()
	env := cfg.ServerConf.Env
	if env == "local" { // if local development, set expire time to 1 year
		expireTime = time.Now().Add(time.Hour * 8760).Unix()
		// expireTime = time.Now().Add(time.Second * 10).Unix() // 10 seconds token to test in local env
	}
	if tokenType == "refreshToken" {
		expireTime = time.Now().Add(time.Hour * 720).Unix() // 30 days for refresh token?
	}

	claims := &UserClaims{
		UserId:    *user.Id,
		Username:  user.Name,
		TokenType: tokenType,
		StandardClaims: jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(*user.Id)),
			ExpiresAt: expireTime,
		}}

	return auth.GetToken(claims)
}

func GetUserTokenResponse(user *User) (map[string]interface{}, error) {
	accessClaims := GenerateUserToken(*user, "accessToken")
	refreshClaims := GenerateUserToken(*user, "refreshToken")

	cfg.LoadEnvVariables()
	secret := cfg.Jwt.Secret
	accessToken, accessTokenErr := accessClaims.SignedString([]byte(secret))
	refreshToken, refreshTokenErr := refreshClaims.SignedString([]byte(secret))
	if accessTokenErr != nil || refreshTokenErr != nil {
		return nil, fmt.Errorf("failed to make jwt: %+v", errors.Join(accessTokenErr, refreshTokenErr).Error())
	}

	return map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         *user,
	}, nil
}

func hashUserPassword(pwd *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*pwd), bcrypt.MinCost)
	if err != nil {
		return err
	}

	*pwd = string(hash)
	return nil
}

func (s *Service) Get(queries map[string]interface{}) ([]*User, *helper.Pagination) {
	fmt.Printf("user service get\n")
	users, pagination := s.repo.Get(queries)

	return users, pagination
}

func (s *Service) GetById(queries map[string]interface{}) ([]*User, error) {
	fmt.Printf("user service getById\n")

	records, _ := s.repo.Get(queries)
	if len(records) == 0 {
		return nil, fmt.Errorf("%s with id: %s not found", tableName, queries["id"])
	}
	return records, nil
}

func (s *Service) Create(users []*User) ([]*User, *helper.HttpErr) {
	fmt.Printf("user service create\n")
	newUserNames := []string{}
	for _, user := range users {
		newUserNames = append(newUserNames, user.Name)
		// hash plain password
		if err := hashUserPassword(user.Password); err != nil {
			return nil, &helper.HttpErr{fiber.StatusInternalServerError, err}
		}
	}

	// check if duplicated by "name"
	existingUsers, _ := s.repo.Get(map[string]interface{}{"name": newUserNames})
	if len(existingUsers) > 0 {
		errMsg := fmt.Sprintf("user service create error: provided user name(s) %+v already exists.\n", newUserNames)
		log.Printf(errMsg)
		return nil, &helper.HttpErr{fiber.StatusConflict, fmt.Errorf(errMsg)}
	}

	results, err := s.repo.Create(users)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Update(users []*User) ([]*User, *helper.HttpErr) {
	fmt.Printf("user service update\n")

	userIds := []string{}
	for _, user := range users {
		userIds = append(userIds, strconv.Itoa(int(*user.Id)))
	}

	// create map by existing user from DB
	userIdMap := map[string]User{}
	existings, _ := s.repo.Get(map[string]interface{}{"id": userIds})
	for _, user := range existings {
		userIdMap[strconv.Itoa(int(*user.Id))] = *user
	}

	// check reqJson for non-existing ids
	// also reuse the map storing the req's user which use for create the "update data"
	nonExistIds := []int64{}
	for _, reqUser := range users {
		_, ok := userIdMap[strconv.Itoa(int(*reqUser.Id))]
		if !ok {
			nonExistIds = append(nonExistIds, *reqUser.Id)
		}
		userIdMap[strconv.Itoa(int(*reqUser.Id))] = *reqUser
	}

	if len(nonExistIds) > 0 || len(existings) == 0 {
		respCode = fiber.StatusNotFound
		notFoundMsg := fmt.Sprintf("cannot update non-existing id(s): %+v", nonExistIds)
		httpErr := &helper.HttpErr{
			Code: fiber.StatusNotFound,
			Err:  errors.New(notFoundMsg),
		}
		return nil, httpErr
	}

	// USELESS, can simply set that column as UNIQUE in DB's table.
	// check conflict of existing name
	for _, user := range users {
		conflicts, _ := s.repo.Get(map[string]interface{}{"name": user.Name})
		if len(conflicts) > 0 && *conflicts[0].Id != *user.Id {
			httpErr := &helper.HttpErr{
				Code: fiber.StatusConflict,
				Err:  fmt.Errorf("%+v is already existed, please try another name.", user.Name),
			}
			return nil, httpErr
		}
	}

	// combining the req user that match with the existing user for update
	for _, originalUser := range existings {
		user := userIdMap[strconv.Itoa(int(*originalUser.Id))] // get the req user
		if user.CreatedAt == nil {
			user.CreatedAt = originalUser.CreatedAt
		}
		if user.Password == nil {
			user.Password = originalUser.Password
		} else {
			hashUserPassword(user.Password)
		}
		newUserBytes, _ := json.Marshal(user)       // convert req user into []byte
		json.Unmarshal(newUserBytes, &originalUser) // unmarshal the req user into its original db record
	}

	results, err := s.repo.Update(existings)
	return results, &helper.HttpErr{fiber.StatusInternalServerError, err}
}

func (s *Service) Delete(ids []string) ([]*User, error) {
	records, _ := s.repo.Get(map[string]interface{}{
		"id": ids,
	})
	if len(records) == 0 {
		return nil, fmt.Errorf("failed to delete, %s with id: %+v not found", tableName, ids)
	}

	return s.repo.Delete(ids)
}

func (s *Service) Login(user *User) (map[string]interface{}, *helper.HttpErr) {
	fmt.Printf("user service login\n")
	results, _ := s.repo.Get(map[string]interface{}{"name": user.Name})
	if len(results) == 0 {
		return nil, &helper.HttpErr{fiber.StatusNotFound, fmt.Errorf("user not exists...")}
	}

	var checkPassword = func(hashedPwd string, plainPwd string) bool {
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
			return false
		}

		return true
	}

	match := checkPassword(*results[0].Password, *user.Password)

	if !match {
		return nil, &helper.HttpErr{fiber.StatusInternalServerError, fmt.Errorf("password not match...")}
	}

	sanitise(results)
	if userTokenResponse, err := GetUserTokenResponse(results[0]); err != nil {
		msg := fmt.Sprintf("failed to refresh token: %+v", err)
		fmt.Println(msg)
		return nil, &helper.HttpErr{fiber.StatusInternalServerError, errors.New(msg)}
	} else {
		return userTokenResponse, nil
	}
}

func (s *Service) Refresh(user *User) (map[string]interface{}, *helper.HttpErr) {
	fmt.Printf("user service login\n")

	results, _ := s.repo.Get(map[string]interface{}{"id": user.Id})
	if len(results) == 0 {
		return nil, &helper.HttpErr{fiber.StatusNotFound, fmt.Errorf("user not exists...")}
	}

	sanitise(results)
	if userTokenResponse, err := GetUserTokenResponse(results[0]); err != nil {
		return nil, &helper.HttpErr{fiber.StatusNotFound, fmt.Errorf("failed to refresh token: %+v", err.Error())}
	} else {
		return userTokenResponse, nil
	}
}
