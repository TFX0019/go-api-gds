package ai

type CreateAIGenerationRequest struct {
	Prompt string `form:"prompt" validate:"required"`
}

type UpdateAIGenerationResultRequest struct {
	ImageOutput  string `form:"image_output" validate:"required_without=ResponseText"`
	ResponseText string `form:"response_text" validate:"required_without=ImageOutput"`
}

type AIGenerationResponse struct {
	ID           string  `json:"id"`
	UserID       uint    `json:"user_id"`
	UserName     string  `json:"user_name"`
	Prompt       string  `json:"prompt"`
	ImageInput   *string `json:"image_input"`
	ImageOutput  *string `json:"image_output"`
	ResponseText *string `json:"response_text"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type PaginatedAIGenerationResponse struct {
	Data  []AIGenerationResponse `json:"data"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
}

type CreateAISuggestionRequest struct {
	Prompt      string `json:"prompt" validate:"required"`
	Description string `json:"description"`
}

type UpdateAISuggestionRequest struct {
	Prompt      string `json:"prompt"`
	Description string `json:"description"`
}

type AISuggestionResponse struct {
	ID          string `json:"id"`
	Prompt      string `json:"prompt"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
