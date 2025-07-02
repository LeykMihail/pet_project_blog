package config

import (
	"os"

	"pet_project_blog/internal/apperrors"
)

// Config содержит базовые настройки
type Config struct {
	Port         string // Порт HTTP сервера
	ConnectBdStr string // Строка подключения к бд
	JWTSecret    string // Секретный ключ для JWT
}

// New создает новую конфигурацию с базовыми настройками или из переменных окружения.
func New() *Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	// postgres://user:pass@localhost:5432/blog?sslmode=disable - пример connStr
	connStr := os.Getenv("DB_CONN_STR")
	secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "s3cr3t_k3y_f0r_jwt_t3st1ng_!@#2025_DEVELOPMENT_ONLY" // Для тестов, заменить на переменную окружения
    }

	return &Config{
		Port:         port,
		ConnectBdStr: connStr,
		JWTSecret: secret,
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
