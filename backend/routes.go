package main

import (
	"os"

	"asterisk-manager/handlers"
	"asterisk-manager/middleware"

	"github.com/gofiber/fiber/v2"
)

func initRoutes(app *fiber.App, h *handlers.Handler, authHandler *handlers.AuthHandler) {
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

	// Auth endpoints (без авторизации)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// Защищенные эндпоинты
	protected := api.Group("/", middleware.JWTAuth(authHandler.GetAuthService()))

	// Auth me endpoint (с авторизацией)
	protected.Get("auth/me", authHandler.Me)

	// Profiles endpoints
	profiles := protected.Group("profiles")
	profiles.Get("/", h.Pagination, h.GetProfiles)
	profiles.Get("/:id", h.GetProfile)
	profiles.Post("/", h.CreateProfile)
	profiles.Put("/:id", h.UpdateProfile)
	profiles.Delete("/:id", h.DeleteProfile)

	// Devices endpoints
	devices := protected.Group("devices")
	devices.Get("/", h.GetDevices)
	devices.Get("/:mac", h.GetDevice)
	devices.Post("/", h.CreateDevice)
	devices.Put("/:mac", h.UpdateDevice)
	devices.Delete("/:mac", h.DeleteDevice)

	// Locations endpoints
	locations := protected.Group("locations")
	locations.Get("/", h.GetLocations)
	locations.Get("/:id", h.GetLocation)
	locations.Post("/", h.CreateLocation)
	locations.Put("/:id", h.UpdateLocation)
	locations.Delete("/:id", h.DeleteLocation)

	// Generator endpoints
	generator := protected.Group("generator")
	_ = generator // TODO: добавить handlers для generator
}
