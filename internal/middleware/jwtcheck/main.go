package jwtcheck

import (
	"golang-api-starter/internal/auth"
	"github.com/gofiber/fiber/v2"
)

/* this is a middleware for checking the jwt */
func CheckFromHeader() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// log.Printf("middleware checking jwt in header.....\n")
		claims, err := auth.ParseJwt(c.Get("Authorization"))
		if err != nil {
			return c.
				Status(fiber.StatusUnauthorized).
				JSON(map[string]interface{}{"message": err.Error()})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}
