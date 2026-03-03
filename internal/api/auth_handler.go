package api

import (
	"encoding/json"
	"net/http"

	"github.com/mati/go-ticket/internal/api/dto"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/mati/go-ticket/internal/services"
)

type AuthHandler struct {
	userService services.UserServiceInterface
}

func NewAuthHandler(userService services.UserServiceInterface) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "User data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	if req.Email == "" || req.Password == "" {
		code, message := MapDomainError(domain.ErrInvalidCredentials)
		ResponseError(w, code, message)
		return
	}

	err := h.userService.RegisterUser(r.Context(), req.Email, req.Password)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	ResponseCreated(w, map[string]string{"message": "User registered successfully"})
}

// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "User data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	if req.Email == "" || req.Password == "" {
		code, message := MapDomainError(domain.ErrInvalidCredentials)
		ResponseError(w, code, message)
		return
	}

	token, err := h.userService.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	ResponseOK(w, map[string]string{"token": token})
}
