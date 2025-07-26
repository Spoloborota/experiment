package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Spoloborota/experiment/internal/domain/repositories"
	"github.com/Spoloborota/experiment/internal/domain/services"
	"github.com/Spoloborota/experiment/internal/interfaces/http/middleware"
)

type ProfileHandler struct {
	profileService *services.ProfileService
	logger         *zap.Logger
}

type CreateProfileRequest struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	Gender    string   `json:"gender"`
	City      string   `json:"city"`
	Interests []string `json:"interests"`
}

type UpdateProfileRequest struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	Gender    string   `json:"gender"`
	City      string   `json:"city"`
	Interests []string `json:"interests"`
}

type ProfilesResponse struct {
	Profiles []interface{} `json:"profiles"`
	Total    int           `json:"total"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
}

func NewProfileHandler(profileService *services.ProfileService, logger *zap.Logger) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         logger,
	}
}

// CreateProfile godoc
// @Summary Создание профиля
// @Description Создает профиль для текущего пользователя
// @Tags profiles
// @Accept json
// @Produce json
// @Param request body CreateProfileRequest true "Данные профиля"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/profile [post]
func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode create profile request", zap.Error(err))
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем профиль
	profile, err := h.profileService.CreateProfile(
		r.Context(),
		user.UserID,
		req.FirstName,
		req.LastName,
		req.Age,
		req.Gender,
		req.City,
		req.Interests,
	)
	if err != nil {
		h.logger.Error("Failed to create profile", zap.Error(err))
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

// GetProfile godoc
// @Summary Получение профиля по ID
// @Description Возвращает профиль пользователя по его ID
// @Tags profiles
// @Produce json
// @Param id path int true "ID профиля"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/profile/{id} [get]
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.writeErrorResponse(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get profile", zap.Error(err))
		h.writeErrorResponse(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// GetMyProfile godoc
// @Summary Получение собственного профиля
// @Description Возвращает профиль текущего авторизованного пользователя
// @Tags profiles
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/profile/me [get]
func (h *ProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	profile, err := h.profileService.GetProfileByUserID(r.Context(), user.UserID)
	if err != nil {
		h.logger.Error("Failed to get user profile", zap.Error(err))
		h.writeErrorResponse(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile godoc
// @Summary Обновление профиля
// @Description Обновляет профиль текущего пользователя
// @Tags profiles
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "Данные для обновления профиля"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/profile/me [put]
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о пользователе из контекста
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode update profile request", zap.Error(err))
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновляем профиль
	profile, err := h.profileService.UpdateProfile(
		r.Context(),
		user.UserID,
		req.FirstName,
		req.LastName,
		req.Age,
		req.Gender,
		req.City,
		req.Interests,
	)
	if err != nil {
		h.logger.Error("Failed to update profile", zap.Error(err))
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// SearchProfiles godoc
// @Summary Поиск профилей
// @Description Ищет профили по заданным фильтрам
// @Tags profiles
// @Produce json
// @Param gender query string false "Фильтр по полу"
// @Param city query string false "Фильтр по городу"
// @Param interests query string false "Фильтр по интересам (через запятую)"
// @Param limit query int false "Лимит результатов" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} ProfilesResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/profiles [get]
func (h *ProfileHandler) SearchProfiles(w http.ResponseWriter, r *http.Request) {
	// Парсим параметры запроса
	filters := repositories.SearchFilters{
		Limit:  10, // По умолчанию
		Offset: 0,
	}

	if gender := r.URL.Query().Get("gender"); gender != "" {
		filters.Gender = &gender
	}

	if city := r.URL.Query().Get("city"); city != "" {
		filters.City = &city
	}

	if interestsStr := r.URL.Query().Get("interests"); interestsStr != "" {
		interests := strings.Split(interestsStr, ",")
		for i, interest := range interests {
			interests[i] = strings.TrimSpace(interest)
		}
		filters.Interests = interests
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// Ищем профили
	profiles, total, err := h.profileService.SearchProfiles(r.Context(), filters)
	if err != nil {
		h.logger.Error("Failed to search profiles", zap.Error(err))
		h.writeErrorResponse(w, "Failed to search profiles", http.StatusInternalServerError)
		return
	}

	// Конвертируем в interface{} для JSON
	profilesInterface := make([]interface{}, len(profiles))
	for i, profile := range profiles {
		profilesInterface[i] = profile
	}

	response := ProfilesResponse{
		Profiles: profilesInterface,
		Total:    total,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProfileHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
