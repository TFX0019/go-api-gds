package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, controller *Controller) {
	route := app.Group("/api/auth")

	route.Post("/register", controller.Register)
	route.Post("/login", controller.Login)
	route.Get("/verify", controller.VerifyEmail)
	route.Post("/refresh", controller.RefreshToken)
	route.Post("/forgot-password", controller.ForgotPassword)
	route.Post("/verify-code", controller.VerifyCode)
	route.Post("/reset-password", controller.ResetPassword)
}
