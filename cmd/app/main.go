package main

import (
	"pet_project_blog/internal/config"
	"pet_project_blog/internal/db"
	"pet_project_blog/internal/handlers"
	"pet_project_blog/internal/logger"
	"pet_project_blog/internal/app"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Создаем logger для разработки
	logger, err := logger.New()
	if err != nil {
		panic("failed to initialize logger " + err.Error())
	}
	logger.Info("Logger initialized successfully")
	defer logger.Sync()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config",
			zap.Error(err),
		)
	}
	logger.Info("Configuration loaded successfully")

	// Применение миграций
	if err := db.ApplyMigrations("file://internal/migrations", cfg.ConnectBdStr); err != nil {
		logger.Fatal("Failed to apply migrations", zap.Error(err))
	}
	logger.Info("Migrations applied successfully")

	// Открытие соединения
	blogDB, err := db.NewDB(cfg.ConnectBdStr)
	if err != nil {
		logger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer blogDB.Close() // Закрытие пула соединений
	logger.Info("Connected to database successfully")

	// Инициализация слоев приложения
	AppLayers := app.InitAppLayers(blogDB, logger, cfg)

	// Инициализация Gin
	r := gin.Default()

	// Настройка доверенных прокси (пустой список для локальной разработки)
	if err := r.SetTrustedProxies([]string{}); err != nil {
		logger.Fatal("Failed to set trusted proxies",
			zap.Error(err),
		)
	}

	// Регистрация маршрутов
	handlers.RegisterRoutes(r, cfg, AppLayers.PostHandler, AppLayers.UserHandler, AppLayers.UserService)

	// Запуск сервера
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
