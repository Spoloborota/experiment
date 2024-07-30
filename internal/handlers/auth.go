package handlers

import (
	"net/http"
	"social_network/internal/auth"
	"social_network/internal/logger"
	"social_network/internal/metrics"
	"social_network/internal/models"
	"social_network/internal/repository"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo  repository.UserRepository
	Validator *validator.Validate
}

func NewAuthHandler(userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo:  userRepo,
		Validator: validator.New(),
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	ctx, span := otel.Tracer("social_network_service").Start(c.Request().Context(), "Register")
	defer span.End()

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		logger.Log.Error("Failed to bind user", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		logger.Log.Error("Failed to hash password", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	user.PasswordHash = string(hashedPassword)

	if err := h.UserRepo.CreateUser(ctx, user); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		logger.Log.Error("Failed to create user", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	metrics.RegisterCounter.Inc()
	return c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx, span := otel.Tracer("social_network_service").Start(c.Request().Context(), "Login")
	defer span.End()

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		logger.Log.Error("Failed to bind user", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	storedUser, err := h.UserRepo.GetUserByName(ctx, user.FirstName, user.LastName)
	if err != nil {
		span.SetAttributes(attribute.String("error", "Invalid credentials"))
		logger.Log.Error("Invalid credentials", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, "Invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		span.SetAttributes(attribute.String("error", "Invalid credentials"))
		logger.Log.Error("Invalid credentials", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := auth.GenerateJWT(storedUser.ID)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		logger.Log.Error("Failed to generate JWT", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	metrics.LoginCounter.Inc()
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *AuthHandler) Validate(i interface{}) error {
	if err := h.Validator.Struct(i); err != nil {
		logger.Log.Error("Validation failed", zap.Error(err))
		return err
	}
	return nil
}
