package products

import (
	"fmt"

	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	var req CreateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	if err := c.service.Create(userID, req); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, nil, "product created successfully")
}

func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	var query PaginationQuery
	if err := ctx.QueryParser(&query); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid query parameters")
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	res, err := c.service.GetAll(query.Page, query.Limit)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "products retrieved successfully")
}

func (c *Controller) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	res, err := c.service.GetByID(id)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusNotFound, "product not found")
	}

	return utils.SendSuccess(ctx, res, "product retrieved successfully")
}

func (c *Controller) GetByUserID(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	if userID == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "user id required")
	}

	var query PaginationQuery
	if err := ctx.QueryParser(&query); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid query parameters")
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	res, err := c.service.GetByUserID(userID, query.Page, query.Limit)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "user products retrieved successfully")
}

func (c *Controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	var req UpdateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.service.Update(id, req); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "product updated successfully")
}

func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	if err := c.service.Delete(id); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "product deleted successfully")
}

func getUserIDFromToken(ctx *fiber.Ctx) (string, error) {
	// Reusing logic from customers/controller ideally, or refactor to utils.
	// For now, duplicating to keep features decoupled.
	userToken := ctx.Locals("user")
	if userToken == nil {
		return "", fmt.Errorf("no user in context")
	}

	token, ok := userToken.(*jwt.Token)
	if !ok {
		return "", fmt.Errorf("invalid token type")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	switch v := claims["user_id"].(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%.0f", v), nil
	default:
		return "", fmt.Errorf("invalid user_id type in token")
	}
}
