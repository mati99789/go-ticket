package dto

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" log:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" log:"-"`
}
