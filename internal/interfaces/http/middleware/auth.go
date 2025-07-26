package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Spoloborota/experiment/internal/domain/services"
)

type contextKey string

const UserContextKey contextKey = "user"

// JWTAuthMiddleware создает middleware для проверки JWT токенов
func JWTAuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Проверяем формат Bearer token
			tokenParts := strings.SplitN(authHeader, " ", 2)
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := tokenParts[1]

			// Валидируем токен
			claims, err := authService.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем информацию о пользователе в контекст
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext извлекает информацию о пользователе из контекста
func GetUserFromContext(ctx context.Context) (*services.JWTClaims, bool) {
	user, ok := ctx.Value(UserContextKey).(*services.JWTClaims)
	return user, ok
}
