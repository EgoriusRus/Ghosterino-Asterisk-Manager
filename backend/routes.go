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

	// Profiles endpoints
	profiles := api.Group("/profiles")
	profiles.Get("/", h.GetProfiles)
	profiles.Get("/:id", h.GetProfile)
	profiles.Post("/", h.CreateProfile)
	profiles.Put("/:id", h.UpdateProfile)
	profiles.Delete("/:id", h.DeleteProfile)

	// Devices endpoints
	devices := api.Group("/devices")
	devices.Get("/", h.GetDevices)
	devices.Get("/:mac", h.GetDevice)
	devices.Post("/", h.CreateDevice)
	devices.Put("/:mac", h.UpdateDevice)
	devices.Delete("/:mac", h.DeleteDevice)

	// Locations endpoints
	locations := api.Group("/locations")
	locations.Get("/", h.GetLocations)
	locations.Get("/:id", h.GetLocation)
	locations.Post("/", h.CreateLocation)
	locations.Put("/:id", h.UpdateLocation)
	locations.Delete("/:id", h.DeleteLocation)

	// Generator endpoints
	generator := api.Group("/generator")
	_ = generator // TODO: добавить handlers для generator
}
