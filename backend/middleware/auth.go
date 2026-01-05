package middleware

import (
	"strings"

	"asterisk-manager/services"

	"github.com/gofiber/fiber/v2"
)

// JWTAuth middleware для проверки JWT токена
func JWTAuth(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем токен из заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Authorization header required")
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := parts[1]

		// Валидируем токен
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Сохраняем claims в контексте
		c.Locals("user", claims)

		return c.Next()
	}
}
