package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type CustomMiddlewares struct {
	PermissionChecker IPermissionCheck
	JwtChecker        IJwtChecker
	Logger            ILogger
}

type IPermissionCheck interface {
	CheckAccess(string) fiber.Handler
}
type IJwtChecker interface {
	CheckJwt() fiber.Handler
}
type ILogger interface {
	Log() fiber.Handler
}
