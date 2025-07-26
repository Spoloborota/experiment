package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Spoloborota/experiment/internal/domain/services"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *services.AuthService
	logger      *zap.Logger
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewAuthHandler(authService *services.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode register request", zap.Error(err))
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем пользователя
	user, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерируем токен
	token, _, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to login after registration", zap.Error(err))
		h.writeErrorResponse(w, "Registration successful but login failed", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Авторизует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для авторизации"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode login request", zap.Error(err))
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Авторизуемся
	token, user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", zap.Error(err))
		h.writeErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
