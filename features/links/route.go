package links

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, controller *Controller) {
	route := app.Group("/api/links", middleware.Protected())

	// Accessible to all roles
	route.Get("/", controller.GetAll)

	// Admin only routes
	adminRoute := route.Group("/", middleware.RequireRole("admin"))
	adminRoute.Post("/", controller.Create)
	adminRoute.Put("/:id", controller.Update)
	adminRoute.Patch("/:id/activate", controller.Activate)
	adminRoute.Patch("/:id/deactivate", controller.Deactivate)
	adminRoute.Delete("/:id", controller.Delete)
}
