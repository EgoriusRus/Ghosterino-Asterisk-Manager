package handlers

import (
	"asterisk-manager/domain"

	"github.com/gofiber/fiber/v2"
)

// GetProfiles возвращает список всех профилей с пагинацией
func (h *Handler) GetProfiles(c *fiber.Ctx) error {
	// Get pagination from context (set by middleware)
	pagination, ok := c.Locals("pagination").(*domain.PaginationInput)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Pagination not found in context")
	}

	// Get profiles with pagination
	profiles, total, err := h.repos.FindProfilesWithLocations(nil, pagination)
	if err != nil {
		return err
	}

	// Build response with pagination metadata
	paginationResponse := domain.PaginationResponse{
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
	}
	paginationResponse.CalculatePages()

	result := domain.PaginatedResult{
		Data:       profiles,
		Pagination: paginationResponse,
	}

	return c.JSON(result)
}

// GetProfile возвращает один профиль по ID
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	var profile domain.Profile
	if err := h.repos.FindByID(&profile, id); err != nil {
		return err
	}
	return c.JSON(profile)
}

// CreateProfile создает новый профиль
func (h *Handler) CreateProfile(c *fiber.Ctx) error {
	var profile domain.Profile
	if err := c.BodyParser(&profile); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repos.Save(&profile); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(profile)
}

// UpdateProfile обновляет существующий профиль
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	var profile domain.Profile

	// Проверяем существование
	if err := h.repos.FindByID(&profile, id); err != nil {
		return err
	}

	// Парсим новые данные
	if err := c.BodyParser(&profile); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Сохраняем
	if err := h.repos.Save(&profile); err != nil {
		return err
	}

	return c.JSON(profile)
}

// DeleteProfile удаляет профиль
func (h *Handler) DeleteProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	var profile domain.Profile

	// Проверяем существование
	if err := h.repos.FindByID(&profile, id); err != nil {
		return err
	}

	// Удаляем
	if err := h.repos.Delete(&profile); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
