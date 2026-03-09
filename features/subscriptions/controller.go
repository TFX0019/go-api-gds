package subscriptions

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}

func (c *Controller) HandleWebhook(ctx *fiber.Ctx) error {
	var payload RevenueCatWebhook

	if err := ctx.BodyParser(&payload); err != nil {
		log.Printf("[RevenueCat Webhook] Failed to parse payload: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	err := c.service.HandleRevenueCatWebhook(payload)
	if err != nil {
		log.Printf("[RevenueCat Webhook] Error processing event: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Always return 200 OK so RevenueCat knows we processed it, even if we ignored some malformed data inside the webhook
	return ctx.Status(fiber.StatusOK).SendString("OK")
}

func (c *Controller) GetTransactions(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	search := ctx.Query("search", "")

	transactions, err := c.service.ListTransactions(page, limit, search)
	if err != nil {
		log.Printf("[Subscriptions] Error fetching transactions: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(transactions)
}
