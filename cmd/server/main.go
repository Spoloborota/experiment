package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	_ "github.com/Spoloborota/experiment/docs" // Импорт для swagger
	"github.com/Spoloborota/experiment/internal/config"
	"github.com/Spoloborota/experiment/internal/domain/services"
	"github.com/Spoloborota/experiment/internal/infrastructure/database"
	"github.com/Spoloborota/experiment/internal/infrastructure/repository"
	"github.com/Spoloborota/experiment/internal/interfaces/http/routes"
)

// @title Social Network API
// @version 1.0
// @description API для социальной сети с регистрацией и просмотром анкет
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите 'Bearer {токен}' для авторизации
func main() {
	// Инициализируем логгер
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	logger.Info("Starting social network server",
		zap.String("port", cfg.Server.Port),
		zap.String("host", cfg.Server.Host))

	// Подключаемся к базе данных
	db, err := database.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("Connected to database successfully")

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)

	// Инициализируем сервисы
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	profileService := services.NewProfileService(profileRepo)

	// Настраиваем роуты
	router := routes.NewRoutes(authService, profileService, logger)
	handler := router.Setup()

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		logger.Info("Server starting", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started successfully")
	logger.Info("API documentation available at: http://" + cfg.Server.Host + ":" + cfg.Server.Port + "/swagger/")

	// Ждем сигнал для остановки
	<-done
	logger.Info("Server is shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
