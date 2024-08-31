package auth

import (
	"golang-api-starter/internal/config"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)


func ParseJwt(token string) (jwt.MapClaims, error) {
	tokenStr := strings.Split(token, "Bearer ")
	// fmt.Println("tokenStr:", len(tokenStr), tokenStr)

	if len(tokenStr) != 2 {
		return nil, logger.Errorf("Malformed token")
	}

	tokenString := tokenStr[1]
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, logger.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secret := config.Cfg.Jwt.Secret
		return []byte(secret), nil
	})

	if err != nil {
		return nil, logger.Errorf(err.Error())
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)

	// fmt.Printf("?? %+v\n", jwtToken)
	// fmt.Println("exp: ", claims["exp"])
	if int64(claims["exp"].(float64)) < time.Now().Local().Unix() {
		err := logger.Errorf("token expired")
		return claims, err
	}

	if !ok && !jwtToken.Valid {
		err := logger.Errorf("Unauthorized")
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
