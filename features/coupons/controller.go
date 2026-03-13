package coupons

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

func (c *Controller) Get(ctx *fiber.Ctx) error {
	res, err := c.service.Get()
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(ctx, res, "coupon retrieved successfully")
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	var req UpdateCouponRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	res, err := c.service.Update(req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "coupon updated successfully")
}

func (c *Controller) Activate(ctx *fiber.Ctx) error {
	if err := c.service.Activate(); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(ctx, nil, "coupon activated successfully")
}

func (c *Controller) Deactivate(ctx *fiber.Ctx) error {
	if err := c.service.Deactivate(); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(ctx, nil, "coupon deactivated successfully")
}
