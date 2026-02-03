package auth

import (
	"github.com/TFX0019/api-go-gds/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, controller *Controller) {
	route := app.Group("/api/auth")

	route.Post("/register", controller.Register)
	route.Post("/login", controller.Login)
	route.Get("/verify", controller.VerifyEmail) // Deprecated?
	route.Post("/verify-account", controller.VerifyAccount)
	route.Post("/resend-code", controller.ResendVerificationCode)
	route.Post("/refresh", controller.RefreshToken)
	route.Post("/forgot-password", controller.ForgotPassword)
	route.Post("/verify-code", controller.VerifyCode)
	route.Post("/reset-password", controller.ResetPassword)
	route.Post("/test-email", controller.TestResendEmail)
	route.Patch("/avatar", middleware.Protected(), controller.UpdateAvatar)
	route.Patch("/name", middleware.Protected(), controller.UpdateName)
}
