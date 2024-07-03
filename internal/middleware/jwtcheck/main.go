package jwtcheck

import (
	"errors"
	"fmt"
	"golang-api-starter/internal/auth"
	"golang-api-starter/internal/config"
	"golang-api-starter/internal/helper"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/helper/utils"
	"golang-api-starter/internal/modules/groupUser"
	"golang-api-starter/internal/modules/user"
	"slices"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type JwtChecker struct{}

var cfg = config.Cfg

/*
* CheckJwt is a middleware for checking the jwt in both cookie & header
* it will first check the cookie, if failed then check the header
 */
func (jc *JwtChecker) CheckJwt(ignorePaths ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// logger.Debugf("middleware checking jwt in header.....")
		url := string(c.Request().URI().Path())
		// logger.Debugf("jwt check: %+v, %+v", url, slices.Contains(ignorePaths, url))
		if slices.Contains(ignorePaths, url) {
			return c.Next()
		}

		requestHeader := c.GetReqHeaders()
		isHtml := strings.Contains(requestHeader["Accept"][0], "text/html")

		var (
			claims jwt.MapClaims
			errStr []string
		)

		var checkUserDisabled = func() error {
			var userId string
			if cfg.DbConf.Driver == "mongodb" {
				userId = claims["userId"].(string)
			} else {
				userId = strconv.Itoa(int(claims["userId"].(float64)))
			}
			err := user.Srvc.IsDisabled(userId)
			if err != nil {
				return err
			}
			return nil
		}

		claims, err := GetTokenFromCookie(c, "accessToken")
		if err != nil {
			if je, ok := err.(*jwtError); ok && je.errorType == "invalid" {
				if claims, err := GetTokenFromCookie(c, "refreshToken"); err != nil {
					errStr = append(errStr, err.Error())
				} else {
					var (
						result     = map[string]interface{}{}
						refreshErr *helper.HttpErr
					)
					if cfg.DbConf.Driver == "mongodb" {
						userId := claims["userId"].(string)
						result, refreshErr = user.Srvc.Refresh(&groupUser.User{MongoId: &userId})
					} else {
						userId := int64(claims["userId"].(float64))
						result, refreshErr = user.Srvc.Refresh(&groupUser.User{Id: utils.ToPtr(helper.FlexInt(userId))})
					}
					if refreshErr != nil {
						errStr = append(errStr, refreshErr.Error())
					}
					if err := user.SetTokensInCookie(result, c); err != nil {
						errStr = append(errStr, err.Error())
					} else {
						logger.Infof(">>>token refreshed")
						c.Locals("claims", claims)
						return c.Next()
					}
				}
			}
			errStr = append(errStr, err.Error())
		} else if userErr := checkUserDisabled(); userErr != nil {
			errStr = append(errStr, userErr.Error())
		} else {
			c.Locals("claims", claims)
			return c.Next()
		}

		claims, err = GetTokenFromHeader(c)
		if err != nil {
			errStr = append(errStr, err.Error())
		} else if userErr := checkUserDisabled(); userErr != nil {
			errStr = append(errStr, userErr.Error())
		}

		if err != nil {
			if isHtml {
				c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
				return c.Redirect("/error", fiber.StatusTemporaryRedirect)
			}

			return c.
				Status(fiber.StatusUnauthorized).
				JSON(map[string]interface{}{"message": errors.Join(errors.New(strings.Join(errStr, ". ")), errors.New("failed to get the jwt from both cookie & header")).Error()})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

func GetTokenFromHeader(ctx *fiber.Ctx) (jwt.MapClaims, error) {
	accessToken := ctx.Get("Authorization")
	if len(accessToken) == 0 {
		return nil, logger.Errorf("Authorization isn't present in header")
	}

	claims, err := auth.ParseJwt(accessToken)
	if err != nil {
		return nil, logger.Errorf("failed to parse token from header, err: %+v", err)
	}

	return claims, nil
}

func GetTokenFromCookie(ctx *fiber.Ctx, tokenType string) (jwt.MapClaims, error) {
	jwt := ctx.Cookies(tokenType)
	if len(jwt) == 0 {
		return nil, &jwtError{
			errorType:    "missing",
			errorMessage: fmt.Sprintf("cookie['%s'] isn't present", tokenType),
		}
	}

	token := "Bearer " + jwt
	// logger.Debugf("%s\n", token)

	claims, err := auth.ParseJwt(token)
	if err != nil {
		return nil, &jwtError{
			errorType:    "invalid",
			errorMessage: fmt.Sprintf("auth.ParseJwt failed, err: %s", err.Error()),
		}
	}

	return claims, nil
}

type jwtError struct {
	errorType    string
	errorMessage string
}

func (je *jwtError) Error() string {
	return je.errorMessage
}
