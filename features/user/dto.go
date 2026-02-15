package user

type UserListDTO struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	IsActive           bool   `json:"is_active"`
	IsVerified         bool   `json:"is_verified"`
	PlanName           string `json:"plan_name"`
	SubscriptionStatus string `json:"subscription_status"`
}
