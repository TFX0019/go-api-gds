package plans

type PlanUpdateDTO struct {
	Title        *string  `json:"title"`
	Description  *string  `json:"description"`
	Price        *float64 `json:"price"`
	Benefits     []string `json:"benefits"`
	MaxCustomers *int     `json:"max_customers"`
	MaxProducts  *int     `json:"max_products"`
	MaxMaterials *int     `json:"max_materials"`
	MaxTasks     *int     `json:"max_tasks"`
	IsActive     *bool    `json:"is_active"`
}

type PlanCreateDTO struct {
	ProductID    string   `json:"product_id" validate:"required"`
	Title        string   `json:"title" validate:"required"`
	Description  string   `json:"description" validate:"required"`
	Price        float64  `json:"price" validate:"required,min=0"`
	Benefits     []string `json:"benefits"`
	MaxCustomers *int     `json:"max_customers" validate:"required"`
	MaxProducts  *int     `json:"max_products" validate:"required"`
	MaxMaterials *int     `json:"max_materials" validate:"required"`
	MaxTasks     *int     `json:"max_tasks" validate:"required"`
	IsActive     *bool    `json:"is_active" default:"true"`
}
