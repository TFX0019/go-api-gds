package banners

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, controller *Controller) {
	route := app.Group("/api/banners", middleware.Protected())

	// Admin only routes
	adminRoute := route.Group("/admin", middleware.RequireRole("admin"))
	adminRoute.Post("/", controller.Create)
	adminRoute.Get("/", controller.GetAllAdmin)
	adminRoute.Delete("/:id", controller.Delete)
}
