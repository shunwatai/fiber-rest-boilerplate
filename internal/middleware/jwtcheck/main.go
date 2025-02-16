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

// CheckJwt is a middleware for checking the jwt in both cookie & header
// it will first check the cookie, if failed then check the header
func (jc *JwtChecker) CheckJwt(ignorePaths ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// logger.Debugf("middleware checking jwt in header.....")
		url := string(c.Request().URI().Path())
		// logger.Debugf("jwt check: %+v, %+v", url, slices.Contains(ignorePaths, url))
		if slices.Contains(ignorePaths, url) {
			return c.Next()
		}

		requestHeader := c.GetReqHeaders()
		if requestHeader["Accept"] == nil || len(strings.TrimSpace(requestHeader["Accept"][0])) == 0 {
			return logger.Errorf("ERROR: missing Accept in request header...")
		}
		isHtml := strings.Contains(requestHeader["Accept"][0], "text/html")

		var (
			claims jwt.MapClaims
			errStr []string
			err    error
		)

		var checkUserDisabled = func() error {
			var userId string
			switch cfg.DbConf.Driver {
			case "mongodb":
				userId = claims["userId"].(string)
			default:
				userId = strconv.Itoa(int(claims["userId"].(float64)))
			}
			err := user.Srvc.IsDisabled(userId)
			if err != nil {
				return err
			}
			return nil
		}

		claims, err = GetTokenFromCookie(c, "accessToken")
		if err != nil {
			refreshClaims, err := refreshTokenIfNeeded(c, claims, err)
			if err != nil {
				errStr = append(errStr, err.Error())
			} else {
				logger.Debugf(">>>refreshed token")
				c.Locals("claims", refreshClaims)
				return c.Next()
			}
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

// GetTokenFromHeader parses the accessToken from request's Authorization header
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

// GetTokenFromCookie parses the accessToken from request's cookies
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

func refreshTokenIfNeeded(c *fiber.Ctx, claims jwt.Claims, err error) (jwt.Claims, error) {
	if err == nil {
		return claims, nil
	}

	je, ok := err.(*jwtError)
	if !ok || je.errorType != "invalid" {
		return nil, err
	}

	refreshTokenClaims, err := GetTokenFromCookie(c, "refreshToken")
	if err != nil {
		return nil, err
	}

	var refreshErr *helper.HttpErr
	var result map[string]interface{}
	switch cfg.DbConf.Driver {
	case "mongodb":
		userId := refreshTokenClaims["userId"].(string)
		result, refreshErr = user.Srvc.Refresh(&groupUser.User{MongoId: &userId})
	default:
		userId := int64(refreshTokenClaims["userId"].(float64))
		result, refreshErr = user.Srvc.Refresh(&groupUser.User{Id: utils.ToPtr(helper.FlexInt(userId))})
	}
	if refreshErr != nil {
		return nil, refreshErr
	}

	return refreshTokenClaims, user.SetTokensInCookie(result, c)
}

type jwtError struct {
	errorType    string
	errorMessage string
}

func (je *jwtError) Error() string {
	return je.errorMessage
}
