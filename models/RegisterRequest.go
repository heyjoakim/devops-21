package models

// RegisterRequest represents a register request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"pwd"`
}
