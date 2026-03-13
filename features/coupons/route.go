package coupons

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, controller *Controller) {
	route := app.Group("/api/coupons", middleware.Protected())

	route.Get("/", controller.Get)

	adminRoute := route.Group("/", middleware.RequireRole("admin"))
	adminRoute.Put("/", controller.Update)
	adminRoute.Patch("/activate", controller.Activate)
	adminRoute.Patch("/deactivate", controller.Deactivate)
}
