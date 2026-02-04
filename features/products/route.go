package products

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, controller *Controller) {
	route := app.Group("/api/products", middleware.Protected())

	route.Post("/", controller.Create)
	route.Get("/", controller.GetAll)
	route.Get("/user", controller.GetByUserID)
	route.Get("/profit-loss", controller.GetProfitLoss)
	route.Get("/:id", controller.GetByID)
	route.Put("/:id", controller.Update)
	route.Patch("/:id/status", controller.UpdateStatus)
	route.Delete("/:id", controller.Delete)
}
