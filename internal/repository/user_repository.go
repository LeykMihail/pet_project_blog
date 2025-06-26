package repository

import (
	"context"
	"database/sql"
	"fmt"

	"pet_project_blog/internal/apperrors"
	"pet_project_blog/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface{
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Выполнение SQL запроса для создания нового пользователя с возвратом ID
	var id int
	err := ur.db.QueryRowContext(ctx,
		`INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id`,
		user.Email, user.PasswordHash,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}

	// Устанавливаем ID созданного поста
	user.ID = id

	return nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	// Выполняем SQL запрос для получения user по email
	err := ur.db.GetContext(ctx, &user,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrSqlNoFoundRows
		}
		return nil, fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}
	return &user, nil
}