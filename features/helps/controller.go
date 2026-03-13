package helps

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

func (c *Controller) Create(ctx *fiber.Ctx) error {
	var req CreateHelpRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	res, err := c.service.Create(req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, res, "help item created successfully")
}

func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	res, err := c.service.GetAll()
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(ctx, res, "helps retrieved successfully")
}

func (c *Controller) GetByTag(ctx *fiber.Ctx) error {
	tag := ctx.Params("tag")
	if tag == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "tag required")
	}

	res, err := c.service.GetByTag(tag)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(ctx, res, "help retrieved successfully")
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	var req UpdateHelpRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := c.service.Update(id, req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "help updated successfully")
}

func (c *Controller) Activate(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	if err := c.service.Activate(id); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(ctx, nil, "help activated successfully")
}

func (c *Controller) Deactivate(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	if err := c.service.Deactivate(id); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(ctx, nil, "help deactivated successfully")
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	if err := c.service.Delete(id); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(ctx, nil, "help deleted successfully")
}
