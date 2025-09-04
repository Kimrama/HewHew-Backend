package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"hewhew-backend/utils"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing bearer token"})
		}
		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := utils.VerifyJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("claims", claims) // ฝาก claims ให้ controller ใช้
		return c.Next()
	}
}
