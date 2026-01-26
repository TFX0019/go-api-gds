package tasks

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, controller *Controller) {
	route := app.Group("/api/tasks", middleware.Protected())

	route.Post("/", controller.Create)
	route.Get("/", controller.GetAll)
	route.Get("/user", controller.GetByUserID)
	route.Get("/:id", controller.GetByID)
	route.Put("/:id", controller.Update)
	route.Delete("/:id", controller.Delete)
}
