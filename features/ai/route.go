package ai

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, controller *Controller) {
	route := app.Group("/api/ai", middleware.Protected())

	route.Post("/", controller.Create)
	route.Patch("/:id/result", controller.UpdateResult)
	route.Get("/me", controller.GetUserGenerations)

	// AI Suggestions (All roles)
	route.Get("/suggestions", controller.GetAllSuggestions)

	// Admin routes
	adminRoute := route.Group("/admin", middleware.RequireRole("admin"))
	adminRoute.Get("/", controller.GetAllAdmin)

	// Admin AI Suggestions
	adminRoute.Post("/suggestions", controller.CreateSuggestion)
	adminRoute.Put("/suggestions/:id", controller.UpdateSuggestion)
	adminRoute.Delete("/suggestions/:id", controller.DeleteSuggestion)
}
