package db

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// NewDB открывает соединение с базой данных и проверяет его
func NewDB(connectStr string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", connectStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// ApplyMigrations применяет миграции к базе данных
func ApplyMigrations(migrationsPath, connectStr string) error {
	m, err := migrate.New(migrationsPath, connectStr)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
