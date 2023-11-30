package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-api-starter/internal/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var cfg = config.Cfg

func ParseJwt(token string) (jwt.MapClaims, error) {
	tokenStr := strings.Split(token, "Bearer ")
	// fmt.Println("tokenStr:", len(tokenStr), tokenStr)

	if len(tokenStr) != 2 {
		errResp, _ := json.Marshal(map[string]string{"error": "Malformed token"})
		return nil, fmt.Errorf(string(errResp))
	}

	tokenString := tokenStr[1]
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		cfg.LoadEnvVariables()
		secret := cfg.Jwt.Secret
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("ParseJwt error: ", err)
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	// fmt.Printf("?? %+v\n", jwtToken)
	// fmt.Println("exp: ", claims["exp"])
	if int64(claims["exp"].(float64)) < time.Now().Local().Unix() {
		err := errors.New("token expired")
		return claims, err
	}

	if !ok && !jwtToken.Valid {
		err := errors.New("Unauthorized")
		return claims, err
	}

	// Access context values in handlers like this
	// props, _ := r.Context().Value("props").(jwt.MapClaims)
	// fmt.Println("props", props)

	return claims, nil
}




func GetToken(claims jwt.Claims) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token
}
