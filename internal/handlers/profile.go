package handlers

import (
	"net/http"
	"social_network/internal/repository"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type ProfileHandler struct {
	UserRepo repository.UserRepository
}

func NewProfileHandler(userRepo repository.UserRepository) *ProfileHandler {
	return &ProfileHandler{UserRepo: userRepo}
}

func (h *ProfileHandler) ViewProfile(c echo.Context) error {
	ctx, span := otel.Tracer("social_network_service").Start(c.Request().Context(), "ViewProfile")
	defer span.End()

	userID := c.Get("user_id").(int)
	user, err := h.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		span.SetAttributes(attribute.String("error", "User not found"))
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}
