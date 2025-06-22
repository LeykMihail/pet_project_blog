package main

import (
	"pet_project_blog/internal/config"
	"pet_project_blog/internal/handlers"
	"pet_project_blog/internal/repository"
	"pet_project_blog/internal/services"

	"github.com/gin-gonic/gin"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Создаем logger для разработки
	loggerConfig := zap.NewDevelopmentConfig()

    // Убираем stacktrace из WARN
    loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	// Настройка отображения времени
    loggerConfig.EncoderConfig.TimeKey = "timestamp"
    loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    // Stacktrace только с Error и выше
    loggerConfig.DisableStacktrace = true // отключает по умолчанию
    logger, err := loggerConfig.Build(
        zap.AddStacktrace(zap.ErrorLevel), // явно добавляем только с error
    )

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

	// Открытие соединения
	blogDB, err := sqlx.Open("pgx", cfg.ConnectBdStr)
	if err != nil {
		logger.Fatal("Unable to connect to database",
			zap.Error(err),
		)
	}
	defer blogDB.Close() // Закрытие пула соединений

	// Проверка соединения
	err = blogDB.Ping()
	if err != nil {
		logger.Fatal("Ping failed",
			zap.Error(err),
		)
	}
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
	postRepo := repository.NewPostRepository(blogDB)
	postService := services.NewPostService(postRepo, logger)
	postHandler := handlers.NewPostHandler(postService, logger)

	// Инициализация Gin
	r := gin.Default()

	// Настройка доверенных прокси (пустой список для локальной разработки)
	if err := r.SetTrustedProxies([]string{}); err != nil {
		logger.Fatal("Failed to set trusted proxies",
			zap.Error(err),
		)
	}

	// Регистрация маршрутов
	handlers.RegisterRoutes(r, postHandler)

	// Запуск сервера
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to start server",
			zap.Error(err),
		)
	}
}
