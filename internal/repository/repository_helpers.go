package repository

import (
	"errors"

	"pet_project_blog/internal/apperrors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Проверяет является ли оштбка нарушение foreign key, если да, то возвращает ошибку, иначе nil
func checkErrForeignKeyViolation(err error) error {
	var pgErr *pgconn.PgError
	// errors.As проверяет, является ли err ошибкой типа *pgconn.PgError, если да то извлекает err в pgErr
	if errors.As(err, &pgErr) {
		// 23503 - это код ошибки "foreign_key_violation"
		if pgErr.Code == "23503" {
			return apperrors.ErrSqlForignKeyViolation
		}
	}
	return nil
}
