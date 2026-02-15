package plans

import (
	"strconv"

	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	ListAll(c *fiber.Ctx) error
	ListActive(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Activate(c *fiber.Ctx) error
	Deactivate(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
}

type controller struct {
	service  Service
	validate *validator.Validate
}

func NewController(service Service) Controller {
	return &controller{
		service:  service,
		validate: validator.New(),
	}
}

func (ctrl *controller) ListAll(c *fiber.Ctx) error {
	plans, err := ctrl.service.ListAll()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"ok":   true,
		"data": plans,
	})
}

func (ctrl *controller) ListActive(c *fiber.Ctx) error {
	plans, err := ctrl.service.ListActive()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"ok":   true,
		"data": plans,
	})
}

func (ctrl *controller) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid plan id")
	}

	var dto PlanUpdateDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	plan, err := ctrl.service.Update(uint(id), dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"data": plan,
	})
}

func (ctrl *controller) Activate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid plan id")
	}

	plan, err := ctrl.service.Activate(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"data": plan,
	})
}

func (ctrl *controller) Deactivate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid plan id")
	}

	plan, err := ctrl.service.Deactivate(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"ok":   true,
		"data": plan,
	})
}

func (ctrl *controller) Create(c *fiber.Ctx) error {
	var dto PlanCreateDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := ctrl.validate.Struct(dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	plan, err := ctrl.service.Create(dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"ok":   true,
		"data": plan,
	})
}
