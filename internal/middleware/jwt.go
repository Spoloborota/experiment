package middleware

import (
	"net/http"
	"social_network/internal/auth"

	"github.com/labstack/echo/v4"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := c.Request().Header.Get("Authorization")
		if tokenStr == "" {
			return c.JSON(http.StatusUnauthorized, "Missing token")
		}

		claims, err := auth.ParseJWT(tokenStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Invalid token")
		}

		c.Set("user_id", claims.UserID)
		return next(c)
	}
}
