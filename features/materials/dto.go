package materials

type CreateMaterialRequest struct {
	Name  string  `form:"name" validate:"required"`
	Price float64 `form:"price" validate:"required,gte=0"`
	Unit  string  `form:"unit" validate:"required"`
}

type UpdateMaterialRequest struct {
	Name  string  `form:"name"`
	Price float64 `form:"price" validate:"gte=0"`
	Unit  string  `form:"unit"`
}

type PaginationQuery struct {
	Page  int `query:"page" validate:"min=1"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}

type MaterialResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Unit      string  `json:"unit"`
	ImageURL  string  `json:"image_url"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type PaginatedResponse struct {
	Data  []MaterialResponse `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
