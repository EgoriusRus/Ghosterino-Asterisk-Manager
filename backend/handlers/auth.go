package handlers

import (
	"asterisk-manager/domain"
	"asterisk-manager/services"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler хендлер авторизации
type AuthHandler struct {
	*Handler
	authService *services.AuthService
}

// NewAuthHandler создает новый хендлер авторизации
func NewAuthHandler(handler *Handler) *AuthHandler {
	return &AuthHandler{
		Handler:     handler,
		authService: services.NewAuthService(),
	}
}

// LoginRequest запрос на авторизацию
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse ответ с токеном
type LoginResponse struct {
	Token string              `json:"token"`
	User  domain.UserResponse `json:"user"`
}

// Login авторизация пользователя
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Username == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username and password are required")
	}

	// Ищем пользователя
	var user domain.User
	if err := h.repos.FindUserByUsername(&user, req.Username); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Проверяем пароль
	if !user.CheckPassword(req.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Генерируем токен
	token, err := h.authService.GenerateToken(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

// Me возвращает текущего пользователя
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims := c.Locals("user").(*services.JWTClaims)

	var user domain.User
	if err := h.repos.FindByID(&user, claims.UserID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return c.JSON(user.ToResponse())
}

// GetAuthService возвращает сервис авторизации для middleware
func (h *AuthHandler) GetAuthService() *services.AuthService {
	return h.authService
}
