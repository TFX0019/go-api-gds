package subscriptions

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, controller *Controller) {
	// RevenueCat webhook endpoint
	webhookGroup := app.Group("/webhooks")

	// Create the RevenueCat specific route
	webhookGroup.Post("/revenuecat", controller.HandleWebhook)
}
