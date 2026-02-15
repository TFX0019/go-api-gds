package plans

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	controller := NewController(service)

	plans := app.Group("/plans")

	// Public or Authenticated (User) routes
	plans.Get("/active", controller.ListActive) // Public list of active plans

	// Admin routes
	plans.Get("/", middleware.Protected(), middleware.RequireRole("admin"), controller.ListAll)
	plans.Post("/", middleware.Protected(), middleware.RequireRole("admin"), controller.Create)
	plans.Put("/:id", middleware.Protected(), middleware.RequireRole("admin"), controller.Update)
	plans.Patch("/:id/activate", middleware.Protected(), middleware.RequireRole("admin"), controller.Activate)
	plans.Patch("/:id/deactivate", middleware.Protected(), middleware.RequireRole("admin"), controller.Deactivate)
}
