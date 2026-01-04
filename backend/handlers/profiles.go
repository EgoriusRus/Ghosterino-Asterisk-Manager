package handlers

import (
	"romchek-asteriska/domain"

	"github.com/gofiber/fiber/v2"
)

// GetProfiles возвращает список всех профилей
func (h *Handler) GetProfiles(c *fiber.Ctx) error {
	profiles, err := h.repos.FindProfilesWithLocations(nil)
	if err != nil {
		return err
	}
	return c.JSON(profiles)
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
