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

	// Admin routes
	adminRoute := route.Group("/admin", middleware.RequireRole("admin"))
	adminRoute.Get("/", controller.GetAllAdmin)
}
