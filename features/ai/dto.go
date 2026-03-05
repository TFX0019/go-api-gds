package ai

type CreateAIGenerationRequest struct {
	Prompt string `form:"prompt" validate:"required"`
}

type UpdateAIGenerationResultRequest struct {
	ImageOutput string `form:"image_output" validate:"required"`
}

type AIGenerationResponse struct {
	ID          string  `json:"id"`
	UserID      uint    `json:"user_id"`
	UserName    string  `json:"user_name"`
	Prompt      string  `json:"prompt"`
	ImageInput  *string `json:"image_input"`
	ImageOutput *string `json:"image_output"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type PaginatedAIGenerationResponse struct {
	Data  []AIGenerationResponse `json:"data"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
}
