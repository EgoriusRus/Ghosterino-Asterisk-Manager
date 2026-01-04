package handlers

import (
	"romchek-asteriska/domain"

	"github.com/gofiber/fiber/v2"
)

// GetLocations возвращает список всех локаций
func (h *Handler) GetLocations(c *fiber.Ctx) error {
	var locations []domain.Location
	if err := h.repos.FindAll(&locations); err != nil {
		return err
	}
	return c.JSON(locations)
}

// GetLocation возвращает одну локацию по ID
func (h *Handler) GetLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	var location domain.Location
	if err := h.repos.FindByID(&location, id); err != nil {
		return err
	}
	return c.JSON(location)
}

// CreateLocation создает новую локацию
func (h *Handler) CreateLocation(c *fiber.Ctx) error {
	var location domain.Location
	if err := c.BodyParser(&location); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repos.Save(&location); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(location)
}

// UpdateLocation обновляет существующую локацию
func (h *Handler) UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	var location domain.Location

	// Проверяем существование
	if err := h.repos.FindByID(&location, id); err != nil {
		return err
	}

	// Парсим новые данные
	if err := c.BodyParser(&location); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Сохраняем
	if err := h.repos.Save(&location); err != nil {
		return err
	}

	return c.JSON(location)
}

// DeleteLocation удаляет локацию
func (h *Handler) DeleteLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	var location domain.Location

	// Проверяем существование
	if err := h.repos.FindByID(&location, id); err != nil {
		return err
	}

	// Удаляем
	if err := h.repos.Delete(&location); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
