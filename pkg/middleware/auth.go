package middleware

import (
	"strings"

	"github.com/TFX0019/api-go-gds/pkg/config"
	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendError(c, fiber.StatusUnauthorized, "missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.SendError(c, fiber.StatusUnauthorized, "invalid authorization header format")
		}

		tokenString := parts[1]
		secret := config.GetEnv("JWT_ACCESS_SECRET", "access_secret")

		token, err := utils.ValidateToken(tokenString, secret)
		if err != nil || !token.Valid {
			return utils.SendError(c, fiber.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals("user", token)
		return c.Next()
	}
}
