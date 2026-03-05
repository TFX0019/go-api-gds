package support

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, controller *Controller) {
	// Support Categories
	catRoute := app.Group("/api/support-categories", middleware.Protected())
	catRoute.Get("/", controller.GetAllCategories)

	// Admin only category routes
	adminCatRoute := catRoute.Group("/", middleware.RequireRole("admin"))
	adminCatRoute.Post("/", controller.CreateCategory)
	adminCatRoute.Put("/:id", controller.UpdateCategory)
	adminCatRoute.Patch("/:id/deactivate", controller.DeactivateCategory)
	adminCatRoute.Patch("/:id/activate", controller.ActivateCategory)

	// Support Tickets
	supRoute := app.Group("/api/supports", middleware.Protected())
	supRoute.Post("/", controller.CreateSupport)
	supRoute.Get("/user", controller.GetUserSupports)
	supRoute.Get("/:id/replies", controller.GetSupportReplies)
	supRoute.Put("/:id", controller.UpdateSupport)
	supRoute.Patch("/:id/status", controller.ChangeSupportStatus)
	supRoute.Delete("/:id", controller.DeleteSupport)

	// Admin only support routes
	adminSupRoute := supRoute.Group("/admin", middleware.RequireRole("admin"))
	adminSupRoute.Get("/", controller.GetAllSupportsAdmin)
}
