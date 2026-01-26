package auth

import (
	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service  Service
	validate *validator.Validate
}

func NewController(service Service) *Controller {
	return &Controller{
		service:  service,
		validate: validator.New(),
	}
}

func (c *Controller) Register(ctx *fiber.Ctx) error {
	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	if err := c.service.Register(req); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, nil, "registration successful, please verify your email")
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	accessToken, refreshToken, user, err := c.service.Login(req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SendSuccess(ctx, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	}, "login successful")
}

func (c *Controller) VerifyEmail(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	if token == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "token required")
	}

	if err := c.service.VerifyEmail(token); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "email verified successfully")
}

func (c *Controller) RefreshToken(ctx *fiber.Ctx) error {
	var req RefreshRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	newAccess, err := c.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SendSuccess(ctx, fiber.Map{"access_token": newAccess}, "token refreshed")
}

func (c *Controller) ForgotPassword(ctx *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	if err := c.service.ForgotPassword(req); err != nil {
		// Do not leak if user exists or not, but for debugging I return error
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "if email exists, recovery code sent")
}

func (c *Controller) VerifyCode(ctx *fiber.Ctx) error {
	var req VerifyCodeRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	if err := c.service.VerifyCode(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "code verified")
}

func (c *Controller) ResetPassword(ctx *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	if err := c.service.ResetPassword(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "password reset successful")
}
