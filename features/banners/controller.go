package banners

import (
	"fmt"
	"path/filepath"

	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) Create(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("image")
	if err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "image file is required")
	}

	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	path := filepath.Join("uploads", filename)
	if err := ctx.SaveFile(file, path); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, "failed to save banner image")
	}

	imageStr := fmt.Sprintf("uploads/%s", filename)

	res, err := c.service.Create(imageStr)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, res, "banner created successfully")
}

func (c *Controller) GetAllAdmin(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	res, err := c.service.GetAllAdmin(page, limit)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "banners retrieved successfully")
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	if err := c.service.Delete(id); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "banner deleted successfully")
}
