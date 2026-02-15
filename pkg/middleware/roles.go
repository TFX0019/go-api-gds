package middleware

import (
	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from locals (set by Protected middleware)
		userToken := c.Locals("user")
		if userToken == nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized")
		}

		token, ok := userToken.(*jwt.Token)
		if !ok {
			return utils.SendError(c, fiber.StatusInternalServerError, "invalid token format")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.SendError(c, fiber.StatusInternalServerError, "invalid token claims")
		}

		rolesInterface, ok := claims["roles"]
		if !ok {
			return utils.SendError(c, fiber.StatusForbidden, "access denied: no roles found")
		}

		roles := convertRoles(rolesInterface)

		for _, r := range roles {
			if r == role {
				return c.Next()
			}
		}

		return utils.SendError(c, fiber.StatusForbidden, "access denied: insufficient permissions")
	}
}

func convertRoles(rolesInterface interface{}) []string {
	var roles []string
	switch v := rolesInterface.(type) {
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok {
				roles = append(roles, str)
			}
		}
	case []string:
		roles = v
	}
	return roles
}
