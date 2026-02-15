package user

import (
	"strconv"

	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	ListUsers(c *fiber.Ctx) error
	Activate(c *fiber.Ctx) error
	Ban(c *fiber.Ctx) error
}

type controller struct {
	service Service
}

func NewController(service Service) Controller {
	return &controller{service: service}
}

func (ctrl *controller) ListUsers(c *fiber.Ctx) error {
	pagination := utils.GetPaginationFromCtx(c)
	result, err := ctrl.service.ListUsers(pagination)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"ok":   true,
		"data": result, // Result is *utils.Pagination which has Rows, TotalRows, etc.
	})
}

func (ctrl *controller) Activate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid user id")
	}

	user, err := ctrl.service.ActivateUser(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"data": user,
	})
}

func (ctrl *controller) Ban(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid user id")
	}

	user, err := ctrl.service.BanUser(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"data": user,
	})
}
