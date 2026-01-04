package handlers

import (
	"asterisk-manager/domain"

	"github.com/gofiber/fiber/v2"
)

// GetDevices возвращает список всех устройств
func (h *Handler) GetDevices(c *fiber.Ctx) error {
	var devices []domain.Device
	if err := h.repos.FindAll(&devices); err != nil {
		return err
	}
	return c.JSON(devices)
}

// GetDevice возвращает одно устройство по MAC
func (h *Handler) GetDevice(c *fiber.Ctx) error {
	mac := c.Params("mac")
	var device domain.Device
	if err := h.repos.FindOne(&device, "mac = ?", mac); err != nil {
		return err
	}
	return c.JSON(device)
}

// CreateDevice создает новое устройство
func (h *Handler) CreateDevice(c *fiber.Ctx) error {
	var device domain.Device
	if err := c.BodyParser(&device); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.repos.Save(&device); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(device)
}

// UpdateDevice обновляет существующее устройство
func (h *Handler) UpdateDevice(c *fiber.Ctx) error {
	mac := c.Params("mac")
	var device domain.Device

	// Проверяем существование
	if err := h.repos.FindOne(&device, "mac = ?", mac); err != nil {
		return err
	}

	// Парсим новые данные
	if err := c.BodyParser(&device); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Сохраняем
	if err := h.repos.Save(&device); err != nil {
		return err
	}

	return c.JSON(device)
}

// DeleteDevice удаляет устройство
func (h *Handler) DeleteDevice(c *fiber.Ctx) error {
	mac := c.Params("mac")
	var device domain.Device

	// Проверяем существование
	if err := h.repos.FindOne(&device, "mac = ?", mac); err != nil {
		return err
	}

	// Удаляем
	if err := h.repos.Delete(&device); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
