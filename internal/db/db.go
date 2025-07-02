package db

import (
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
