package auth

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Code            string `json:"code" validate:"required,len=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type UpdateNameRequest struct {
	Name string `json:"name" validate:"required"`
}

type UserResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Avatar       *string `json:"avatar"`
	IsPro        bool    `json:"is_pro"`
	Plan         string  `json:"plan"`
	MaxCustomers int     `json:"max_customers"`
	MaxProducts  int     `json:"max_products"`
	MaxMaterials int     `json:"max_materials"`
	MaxTasks     int     `json:"max_tasks"`
}
