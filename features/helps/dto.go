package helps

type CreateHelpRequest struct {
	Tag         string `json:"tag" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateHelpRequest struct {
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

type HelpResponse struct {
	ID          uint   `json:"id"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
