package main

import (
	"social_network/config"
	"social_network/internal/db"
	"social_network/internal/handlers"
	"social_network/internal/logger"
	"social_network/internal/metrics"
	"social_network/internal/middleware"
	"social_network/internal/repository"
	"social_network/internal/tracing"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	_ "social_network/docs"
)

func main() {
	logger.Init()
	defer logger.Log.Sync()

	tp, err := tracing.InitTracer()
	if err != nil {
		logger.Log.Fatal("Error initializing tracer", zap.Error(err))
	}
	defer tracing.ShutdownTracer(tp)

	e := echo.New()
	cfg := config.NewConfig()
	db.Migrate(cfg.DB)

	metrics.Init()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(otelecho.Middleware("social_network_service"))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", cfg.DB)
			return next(c)
		}
	})

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	userRepo := repository.NewUserRepository(cfg.DB)

	authHandler := handlers.NewAuthHandler(userRepo)
	profileHandler := handlers.NewProfileHandler(userRepo)

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
	e.GET("/profile", profileHandler.ViewProfile, middleware.JWTAuth)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := e.Start(":8080"); err != nil {
		logger.Log.Fatal("Error starting server", zap.Error(err))
	}
}
