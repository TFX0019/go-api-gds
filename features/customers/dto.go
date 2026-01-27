package customers

type CreateCustomerRequest struct {
	Name             string   `form:"name" validate:"required"`
	Phone            string   `form:"phone"`
	Email            string   `form:"email" validate:"omitempty,email"`
	UsesStandardSize bool     `form:"uses_standard_size"`
	StandardSize     string   `form:"standard_size"`
	Back             *float64 `form:"back"`
	Neck             *float64 `form:"neck"`
	FrontSize        *float64 `form:"front_size"`
	Armhole          *float64 `form:"armhole"`
	BackSize         *float64 `form:"back_size"`
	BustChest        *float64 `form:"bust_chest"`
	Waist            *float64 `form:"waist"`
	Hip              *float64 `form:"hip"`
	RiseHeight       *float64 `form:"rise_height"`
	SkirtLength      *float64 `form:"skirt_length"`
	PantsLength      *float64 `form:"pants_length"`
	KneeWidth        *float64 `form:"knee_width"`
	HemWidth         *float64 `form:"hem_width"`
	SleeveLength     *float64 `form:"sleeve_length"`
	CuffSize         *float64 `form:"cuff_size"`
}

type UpdateCustomerRequest struct {
	Name             string   `form:"name"`
	Phone            string   `form:"phone"`
	Email            string   `form:"email" validate:"omitempty,email"`
	UsesStandardSize bool     `form:"uses_standard_size"`
	StandardSize     string   `form:"standard_size"`
	Back             *float64 `form:"back"`
	Neck             *float64 `form:"neck"`
	FrontSize        *float64 `form:"front_size"`
	Armhole          *float64 `form:"armhole"`
	BackSize         *float64 `form:"back_size"`
	BustChest        *float64 `form:"bust_chest"`
	Waist            *float64 `form:"waist"`
	Hip              *float64 `form:"hip"`
	RiseHeight       *float64 `form:"rise_height"`
	SkirtLength      *float64 `form:"skirt_length"`
	PantsLength      *float64 `form:"pants_length"`
	KneeWidth        *float64 `form:"knee_width"`
	HemWidth         *float64 `form:"hem_width"`
	SleeveLength     *float64 `form:"sleeve_length"`
	CuffSize         *float64 `form:"cuff_size"`
}

type PaginationQuery struct {
	Page  int `query:"page" validate:"min=1"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}

type CustomerResponse struct {
	ID               string   `json:"id"`
	UserID           string   `json:"user_id"`
	AvatarURL        string   `json:"avatar_url"`
	Name             string   `json:"name"`
	Phone            string   `json:"phone"`
	Email            string   `json:"email"`
	UsesStandardSize bool     `json:"uses_standard_size"`
	StandardSize     string   `json:"standard_size"`
	Back             *float64 `json:"back"`
	Neck             *float64 `json:"neck"`
	FrontSize        *float64 `json:"front_size"`
	Armhole          *float64 `json:"armhole"`
	BackSize         *float64 `json:"back_size"`
	BustChest        *float64 `json:"bust_chest"`
	Waist            *float64 `json:"waist"`
	Hip              *float64 `json:"hip"`
	RiseHeight       *float64 `json:"rise_height"`
	SkirtLength      *float64 `json:"skirt_length"`
	PantsLength      *float64 `json:"pants_length"`
	KneeWidth        *float64 `json:"knee_width"`
	HemWidth         *float64 `json:"hem_width"`
	SleeveLength     *float64 `json:"sleeve_length"`
	CuffSize         *float64 `json:"cuff_size"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

type PaginatedResponse struct {
	Data  []CustomerResponse `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
