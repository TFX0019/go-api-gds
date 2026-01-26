package tasks

type CreateTaskRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Status      string  `json:"status" validate:"required,oneof=pending in_progress completed canceled"`
	DateTime    string  `json:"date_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
	ProductID   *string `json:"product_id" validate:"omitempty,uuid"`
}

type UpdateTaskRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status" validate:"omitempty,oneof=pending in_progress completed canceled"`
	DateTime    string  `json:"date_time" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ProductID   *string `json:"product_id" validate:"omitempty,uuid"`
}

type PaginationQuery struct {
	Page      int    `query:"page" validate:"min=1"`
	Limit     int    `query:"limit" validate:"min=1,max=100"`
	Status    string `query:"status"`
	Available string `query:"available"`                                     // Unused? User asked date filter.
	Date      string `query:"date" validate:"omitempty,datetime=2006-01-02"` // Filter by specific date (YYYY-MM-DD)
}

type ClientInfo struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type ProductInfo struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	Client *ClientInfo `json:"client,omitempty"`
}

type TaskResponse struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Status      string       `json:"status"`
	DateTime    string       `json:"date_time"`
	Product     *ProductInfo `json:"product,omitempty"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

type PaginatedResponse struct {
	Data  []TaskResponse `json:"data"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
