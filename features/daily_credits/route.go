package daily_credits

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, controller *Controller) {
	route := app.Group("/api/daily-credits", middleware.Protected())

	// Can be accessed by anyone who is authenticated (or you can restrict to admin if needed)
	// But usually users need to know the config. We'll leave Get open to authenticated users for now.
	route.Get("/", controller.Get)

	// Admin only can update
	adminRoute := route.Group("/", middleware.RequireRole("admin"))
	adminRoute.Patch("/", controller.Update)
	adminRoute.Put("/", controller.Update)
}
