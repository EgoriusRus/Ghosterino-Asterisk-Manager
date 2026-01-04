package main

import (
	"os"

	"romchek-asteriska/handlers"

	"github.com/gofiber/fiber/v2"
)

func initRoutes(app *fiber.App, h *handlers.Handler) {
	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		version := os.Getenv("APP_VERSION")
		if version == "" {
			version = "dev"
		}
		return c.SendString(version)
	})

	// API группа
	api := app.Group("/api")

	// Locations endpoints
	locations := api.Group("/locations")
	_ = locations // TODO: добавить handlers для locations

	// Devices endpoints
	devices := api.Group("/devices")
	_ = devices // TODO: добавить handlers для devices

	// Profiles endpoints
	profiles := api.Group("/profiles")
	_ = profiles // TODO: добавить handlers для profiles

	// Generator endpoints
	generator := api.Group("/generator")
	_ = generator // TODO: добавить handlers для generator
}
