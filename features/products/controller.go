package products

import (
	"fmt"
	"strconv"
	"strings"

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

	productRes, err := c.service.Create(userID, req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, productRes, "product created successfully")
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
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, "unauthorized")
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

func (c *Controller) GetProfitLoss(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	monthStr := ctx.Query("month")
	var month int
	if monthStr != "" {
		month = parseMonth(monthStr)
		if month == 0 {
			// Optional: Return error if invalid month provided?
			// Or just ignore filter?
			// prompt said "filter by month example ?month=jan", implying valid input.
			// I'll ignore invalid input or better, return error.
			// Let's just pass 0 if invalid.
		}
	}

	res, err := c.service.GetProfitLoss(userID, month)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "profit and loss retrieved successfully")
}

func parseMonth(m string) int {
	if m == "" {
		return 0
	}
	m = strings.ToLower(m)
	switch m {
	case "jan", "january", "1", "01":
		return 1
	case "feb", "february", "2", "02":
		return 2
	case "mar", "march", "3", "03":
		return 3
	case "apr", "april", "4", "04":
		return 4
	case "may", "5", "05":
		return 5
	case "jun", "june", "6", "06":
		return 6
	case "jul", "july", "7", "07":
		return 7
	case "aug", "august", "8", "08":
		return 8
	case "sep", "september", "9", "09":
		return 9
	case "oct", "october", "10":
		return 10
	case "nov", "november", "11":
		return 11
	case "dec", "december", "12":
		return 12
	}
	val, err := strconv.Atoi(m)
	if err == nil && val >= 1 && val <= 12 {
		return val
	}
	return 0
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

	productRes, err := c.service.Update(id, req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, productRes, "product updated successfully")
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

func (c *Controller) UpdateStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	var req UpdateProductStatusRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	productRes, err := c.service.UpdateStatus(id, req.Status)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, productRes, "product status updated successfully")
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
