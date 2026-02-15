package user

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	controller := NewController(service)

	users := app.Group("/users")

	// Admin routes
	users.Get("/", middleware.Protected(), middleware.RequireRole("admin"), controller.ListUsers)
	users.Patch("/:id/activate", middleware.Protected(), middleware.RequireRole("admin"), controller.Activate)
	users.Patch("/:id/ban", middleware.Protected(), middleware.RequireRole("admin"), controller.Ban)
}
