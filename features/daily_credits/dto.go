package daily_credits

type UpdateDailyCreditRequest struct {
	Free    int `json:"free" validate:"required,min=0"`
	Premium int `json:"premium" validate:"required,min=0"`
}

type DailyCreditResponse struct {
	ID        uint   `json:"id"`
	Free      int    `json:"free"`
	Premium   int    `json:"premium"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
