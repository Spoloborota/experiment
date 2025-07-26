package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/Spoloborota/experiment/internal/domain/services"
	"github.com/Spoloborota/experiment/internal/interfaces/http/handlers"
	authMiddleware "github.com/Spoloborota/experiment/internal/interfaces/http/middleware"
)

type Routes struct {
	authService    *services.AuthService
	profileService *services.ProfileService
	logger         *zap.Logger
}

func NewRoutes(authService *services.AuthService, profileService *services.ProfileService, logger *zap.Logger) *Routes {
	return &Routes{
		authService:    authService,
		profileService: profileService,
		logger:         logger,
	}
}

func (rt *Routes) Setup() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler())

	// Создаем обработчики
	authHandler := handlers.NewAuthHandler(rt.authService, rt.logger)
	profileHandler := handlers.NewProfileHandler(rt.profileService, rt.logger)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Публичные роуты (без авторизации)
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Get("/profile/{id}", profileHandler.GetProfile)
		r.Get("/profiles", profileHandler.SearchProfiles)

		// Защищенные роуты (с авторизацией)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.JWTAuthMiddleware(rt.authService))

			r.Get("/profile/me", profileHandler.GetMyProfile)
			r.Post("/profile", profileHandler.CreateProfile)
			r.Put("/profile/me", profileHandler.UpdateProfile)
		})
	})

	return r
}
