package utils

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Ok      bool        `json:"ok"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
}

func SendSuccess(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Ok:      true,
		Data:    data,
		Message: message,
		Error:   "",
	})
}

func SendCreated(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Ok:      true,
		Data:    data,
		Message: message,
		Error:   "",
	})
}

func SendError(c *fiber.Ctx, status int, errStr string) error {
	return c.Status(status).JSON(APIResponse{
		Ok:      false,
		Data:    nil,
		Message: "", // Or maybe we can put the error message here too? The user asked for "error: si hay error".
		Error:   errStr,
	})
}
