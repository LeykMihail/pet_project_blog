package main

import (
	"pet_project_blog/internal/config"
	"pet_project_blog/internal/handlers"
	"pet_project_blog/internal/logger"
	"pet_project_blog/internal/repository"
	"pet_project_blog/internal/services"
	"pet_project_blog/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	// Подключение db
	blogDB, err := db.NewDB(cfg.ConnectBdStr)
	if err != nil {
		logger.Fatal("Unable to connect to database or ping failed", zap.Error(err))
	}
	defer blogDB.Close()
	logger.Info("Connected to database successfully")

	// Настройка миграций
	m, err := migrate.New(
		"file://internal/migrations",
		cfg.ConnectBdStr,
	)
	if err != nil {
		logger.Fatal("Failed to initialize migrations",
			zap.Error(err),
		)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("Failed to apply migrations",
			zap.Error(err),
		)
	}
	logger.Info("Migrations applied successfully")

	// Инициализация слоев приложения
	userRepo := repository.NewUserRepository(blogDB)
	postRepo := repository.NewPostRepository(blogDB)
	userService := services.NewUserService(userRepo, logger)
	postService := services.NewPostService(postRepo, logger)
	postHandler := handlers.NewPostHandler(postService, logger)
	userHandler := handlers.NewUserHandler(userService, logger, cfg)

	// Инициализация Gin
	r := gin.Default()

	// Настройка доверенных прокси (пустой список для локальной разработки)
	if err := r.SetTrustedProxies([]string{}); err != nil {
		logger.Fatal("Failed to set trusted proxies",
			zap.Error(err),
		)
	}

	// Регистрация маршрутов
	handlers.RegisterRoutesPost(r, postHandler, cfg, userService)
	handlers.RegisterRoutesUser(r, userHandler)

	// Запуск сервера
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to start server",
			zap.Error(err),
		)
	}
}
