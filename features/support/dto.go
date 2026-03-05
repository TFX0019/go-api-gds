package support

type CreateCategoryRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Active      *bool  `json:"active"`
}

type CreateSupportRequest struct {
	Subject           string  `form:"subject" validate:"required"`
	Description       string  `form:"description" validate:"required"`
	SupportCategoryID string  `form:"support_category_id" validate:"required,uuid"`
	ParentID          *string `form:"parent_id" validate:"omitempty,uuid"`
}

type UpdateSupportRequest struct {
	Subject           string  `json:"subject"`
	Description       string  `json:"description"`
	SupportCategoryID *string `json:"support_category_id" validate:"omitempty,uuid"`
	Status            *string `json:"status" validate:"omitempty,oneof=open in_process closed"`
}

type PaginationQuery struct {
	Page  int `query:"page" validate:"min=1"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}

type SupportCategoryResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type SupportResponse struct {
	ID                string                   `json:"id"`
	Subject           string                   `json:"subject"`
	Description       string                   `json:"description"`
	UserID            uint                     `json:"user_id"`
	SupportCategoryID string                   `json:"support_category_id"`
	SupportCategory   *SupportCategoryResponse `json:"support_category,omitempty"`
	Status            string                   `json:"status"`
	Image             string                   `json:"image"`
	ParentID          *string                  `json:"parent_id"`
	IsDeleted         bool                     `json:"is_deleted"`
	CreatedAt         string                   `json:"created_at"`
	UpdatedAt         string                   `json:"updated_at"`
}

type PaginatedSupportResponse struct {
	Data  []SupportResponse `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}
