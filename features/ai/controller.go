package ai

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

func (c *Controller) Create(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	var req CreateAIGenerationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	var imageInput *string
	file, err := ctx.FormFile("image_input")
	if err == nil {
		filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
		path := filepath.Join("uploads", filename)
		if err := ctx.SaveFile(file, path); err != nil {
			return utils.SendError(ctx, fiber.StatusInternalServerError, "failed to save input image")
		}
		imgStr := fmt.Sprintf("uploads/%s", filename)
		imageInput = &imgStr
	}

	res, err := c.service.CreateGeneration(userID, req.Prompt, imageInput)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, res, "ai generation record created successfully")
}

func (c *Controller) UpdateResult(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	var imageOutput string
	file, err := ctx.FormFile("image_output")
	if err == nil {
		filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
		path := filepath.Join("uploads", filename)
		if err := ctx.SaveFile(file, path); err != nil {
			return utils.SendError(ctx, fiber.StatusInternalServerError, "failed to save output image")
		}
		imageOutput = fmt.Sprintf("uploads/%s", filename)
	} else {
		// Try to get from body if not a file
		var req UpdateAIGenerationResultRequest
		if err := ctx.BodyParser(&req); err != nil || req.ImageOutput == "" {
			return utils.SendError(ctx, fiber.StatusBadRequest, "image_output is required (file or path)")
		}
		imageOutput = req.ImageOutput
	}

	res, err := c.service.UpdateGenerationResult(id, userID, imageOutput)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "ai generation result saved successfully")
}

func (c *Controller) GetUserGenerations(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromToken(ctx)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusUnauthorized, "unauthorized")
	}

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	res, err := c.service.GetGenerationsByUserID(userID, page, limit)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "ai generations retrieved successfully")
}

func (c *Controller) GetAllAdmin(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	res, err := c.service.GetAllGenerationsAdmin(page, limit)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "all ai generations retrieved successfully")
}

// AI Suggestions Endpoints

func (c *Controller) CreateSuggestion(ctx *fiber.Ctx) error {
	var req CreateAISuggestionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	if err := c.validate.Struct(req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, utils.ParseValidationError(err))
	}

	res, err := c.service.CreateSuggestion(req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendCreated(ctx, res, "ai suggestion created successfully")
}

func (c *Controller) UpdateSuggestion(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	var req UpdateAISuggestionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.SendError(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := c.service.UpdateSuggestion(id, req)
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "ai suggestion updated successfully")
}

func (c *Controller) DeleteSuggestion(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utils.SendError(ctx, fiber.StatusBadRequest, "id required")
	}

	if err := c.service.DeleteSuggestion(id); err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, nil, "ai suggestion deleted successfully")
}

func (c *Controller) GetAllSuggestions(ctx *fiber.Ctx) error {
	res, err := c.service.GetAllSuggestions()
	if err != nil {
		return utils.SendError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccess(ctx, res, "ai suggestions retrieved successfully")
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
