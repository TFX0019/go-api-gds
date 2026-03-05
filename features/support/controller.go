package support

import (
	"fmt"
	"path/filepath"

	"github.com/TFX0019/api-go-gds/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

// Category Endpoints

func (c *Controller) CreateCategory(ctx *fiber.Ctx) error {
	var req CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	res, err := c.service.CreateCategory(req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, res, "support category created successfully")
}

func (c *Controller) GetAllCategories(ctx *fiber.Ctx) error {
	res, err := c.service.GetAllCategories()
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(ctx, res, "support categories retrieved successfully")
}

func (c *Controller) UpdateCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	var req UpdateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := c.service.UpdateCategory(id, req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "support category updated successfully")
}

func (c *Controller) DeactivateCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	if err := c.service.DeactivateCategory(id); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "support category deactivated successfully")
}

func (c *Controller) ActivateCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	if err := c.service.ActivateCategory(id); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "support category activated successfully")
}

// Support Endpoints

func (c *Controller) CreateSupport(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	var req CreateSupportRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	var imageURL string
	file, err := ctx.FormFile("image")
	if err == nil {
		filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
		path := filepath.Join("uploads", filename)
		if err := ctx.SaveFile(file, path); err != nil {
			return utils.SendError(ctx, fiber.StatusInternalServerError, "failed to save image")
		}
		imageURL = fmt.Sprintf("uploads/%s", filename)
	}

	res, err := c.service.CreateSupport(uint(userID), req, imageURL)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, res, "support ticket created successfully")
}

func (c *Controller) GetUserSupports(ctx *fiber.Ctx) error {
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

	res, err := c.service.GetParentSupportsByUserID(uint(userID), query.Page, query.Limit)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "user support tickets retrieved successfully")
}

func (c *Controller) GetAllSupportsAdmin(ctx *fiber.Ctx) error {
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

	res, err := c.service.GetAllParentSupportsAdmin(query.Page, query.Limit)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "all support tickets retrieved successfully")
}

func (c *Controller) GetSupportReplies(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	res, err := c.service.GetRepliesByParentID(id)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "support replies retrieved successfully")
}

func (c *Controller) UpdateSupport(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	var req UpdateSupportRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := c.service.UpdateSupport(id, req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "support ticket updated successfully")
}

func (c *Controller) ChangeSupportStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	var req ChangeSupportStatusRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	res, err := c.service.ChangeSupportStatus(id, req.Status)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "support ticket status updated successfully")
}

func (c *Controller) DeleteSupport(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	isAdmin := hasAdminRole(ctx)

	if err := c.service.DeleteSupport(id, uint(userID), isAdmin); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "support ticket deleted successfully")
}

func getUserIDFromToken(ctx *fiber.Ctx) (uint, error) {
	userToken := ctx.Locals("user")
	if userToken == nil {
		return 0, fmt.Errorf("no user in context")
	}

	token, ok := userToken.(*jwt.Token)
	if !ok {
		return 0, fmt.Errorf("invalid token type")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	switch v := claims["user_id"].(type) {
	case float64:
		return uint(v), nil
	default:
		return 0, fmt.Errorf("invalid user_id type in token")
	}
}

func hasAdminRole(ctx *fiber.Ctx) bool {
	userToken := ctx.Locals("user")
	if userToken == nil {
		return false
	}
	token, ok := userToken.(*jwt.Token)
	if !ok {
		return false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	rolesInterface, ok := claims["roles"]
	if !ok {
		return false
	}

	var roles []string
	switch v := rolesInterface.(type) {
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok {
				roles = append(roles, str)
			}
		}
	case []string:
		roles = v
	}

	for _, r := range roles {
		if r == "admin" {
			return true
		}
	}
	return false
}
