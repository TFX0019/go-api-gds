package links

type CreateLinkRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"required,url"`
}

type UpdateLinkRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url" validate:"omitempty,url"`
}

type LinkResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
