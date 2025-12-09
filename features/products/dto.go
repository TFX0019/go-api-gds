package products

type CreateProductRequest struct {
	Name                 string  `json:"name" validate:"required"`
	ClientID             *string `json:"client_id" validate:"omitempty,uuid"`
	MaterialsCost        float64 `json:"materials_cost" validate:"gte=0"`
	HoursCost            float64 `json:"hours_cost" validate:"gte=0"`
	ProfitPercentage     float64 `json:"profit_percentage" validate:"gte=0"`
	IncludeFixedExpenses bool    `json:"include_fixed_expenses"`
	FixedExpenseRate     float64 `json:"fixed_expense_rate" validate:"gte=0"`
	Subtotal             float64 `json:"subtotal" validate:"gte=0"`
	FixedExpensesAmount  float64 `json:"fixed_expenses_amount" validate:"gte=0"`
	BaseTotal            float64 `json:"base_total" validate:"gte=0"`
	ProfitAmount         float64 `json:"profit_amount" validate:"gte=0"`
	Total                float64 `json:"total" validate:"gte=0"`
}

type UpdateProductRequest struct {
	Name                 string  `json:"name"`
	ClientID             *string `json:"client_id" validate:"omitempty,uuid"`
	MaterialsCost        float64 `json:"materials_cost" validate:"gte=0"`
	HoursCost            float64 `json:"hours_cost" validate:"gte=0"`
	ProfitPercentage     float64 `json:"profit_percentage" validate:"gte=0"`
	IncludeFixedExpenses bool    `json:"include_fixed_expenses"`
	FixedExpenseRate     float64 `json:"fixed_expense_rate" validate:"gte=0"`
	Subtotal             float64 `json:"subtotal" validate:"gte=0"`
	FixedExpensesAmount  float64 `json:"fixed_expenses_amount" validate:"gte=0"`
	BaseTotal            float64 `json:"base_total" validate:"gte=0"`
	ProfitAmount         float64 `json:"profit_amount" validate:"gte=0"`
	Total                float64 `json:"total" validate:"gte=0"`
}

type PaginationQuery struct {
	Page  int `query:"page" validate:"min=1"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}

type ProductResponse struct {
	ID                   string  `json:"id"`
	UserID               string  `json:"user_id"`
	Name                 string  `json:"name"`
	ClientID             *string `json:"client_id,omitempty"`
	MaterialsCost        float64 `json:"materials_cost"`
	HoursCost            float64 `json:"hours_cost"`
	ProfitPercentage     float64 `json:"profit_percentage"`
	IncludeFixedExpenses bool    `json:"include_fixed_expenses"`
	FixedExpenseRate     float64 `json:"fixed_expense_rate"`
	Subtotal             float64 `json:"subtotal"`
	FixedExpensesAmount  float64 `json:"fixed_expenses_amount"`
	BaseTotal            float64 `json:"base_total"`
	ProfitAmount         float64 `json:"profit_amount"`
	Total                float64 `json:"total"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

type PaginatedResponse struct {
	Data  []ProductResponse `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}
