package handlers

import (
	"log"

	"asterisk-manager/domain"
	"asterisk-manager/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	repos *repositories.Repos
}

func NewHandler(repos *repositories.Repos) *Handler {
	return &Handler{
		repos: repos,
	}
}

// ErrorResponse структура ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// Pagination middleware parses and validates pagination query parameters
func (h *Handler) Pagination(c *fiber.Ctx) error {
	var pagination domain.PaginationInput

	// Parse query parameters
	if err := c.QueryParser(&pagination); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid pagination parameters")
	}

	// Validate and set defaults
	pagination.Validate()

	// Store in context
	c.Locals("pagination", &pagination)

	return c.Next()
}

// ErrorHandler централизованная обработка ошибок
func (h *Handler) ErrorHandler(ctx *fiber.Ctx, err error) error {
	// GORM Record Not Found
	if err == gorm.ErrRecordNotFound {
		log.Printf("Record not found: %s %s", ctx.Method(), ctx.Path())
		return ctx.Status(404).JSON(ErrorResponse{Error: "Not found"})
	}

	// Fiber ошибки
	if fiberError, ok := err.(*fiber.Error); ok {
		log.Printf("Fiber error: %s %s - %s", ctx.Method(), ctx.Path(), fiberError.Error())
		return ctx.Status(fiberError.Code).JSON(ErrorResponse{Error: fiberError.Message})
	}

	// Все остальные ошибки
	log.Printf("Unexpected error: %s %s - %v", ctx.Method(), ctx.Path(), err)
	return ctx.Status(500).JSON(ErrorResponse{Error: "Internal server error"})
}
