package dto

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" log:"-"` //nolint:gosec // G706: slog uses structured fields
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password" log:"-"` //nolint:gosec // G706: slog uses structured fields
}
