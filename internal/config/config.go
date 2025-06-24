package config

import (
	"os"

	"pet_project_blog/internal/apperrors"
)

// Config содержит базовые настройки
type Config struct {
	Port         string // Порт HTTP сервера
	ConnectBdStr string // Строка подключения к бд
}

// New создает новую конфигурацию с базовыми настройками или из переменных окружения.
func New() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	// postgres://user:pass@localhost:5432/blog?sslmode=disable - пример connStr
	connStr := os.Getenv("DB_CONN_STR")

	return &Config{
		Port:         port,
		ConnectBdStr: connStr,
	}
}

// Load загружает конфигурацию
func Load() (*Config, error) {
	cfg := New()
	if cfg.ConnectBdStr == "" {
        return nil, apperrors.ErrConfigEmptyDB_CONN_STR
    }
    return cfg, nil
}
