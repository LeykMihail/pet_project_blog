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

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)

	CreateSubscription(ctx context.Context, userID, authorID int) error
    GetSubscriptionsByUserID(ctx context.Context, userID int) ([]int, error)
    DeleteSubscription(ctx context.Context, userID, authorID int) error
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
		if uErr := checkErrUniqueViolation(err); uErr != nil {
			return uErr
		}
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

func (ur *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	// Выполняем SQL запрос для получения user по id
	err := ur.db.GetContext(ctx, &user,
		`SELECT id, email, password_hash, created_at FROM users WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrSqlNoFoundRows
		}
		return nil, fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}
	return &user, nil
}

func (ur *userRepository) CreateSubscription(ctx context.Context, userID, authorID int) error {
    _, err := ur.db.ExecContext(ctx,
        `INSERT INTO subscriptions (user_id, author_id) VALUES ($1, $2) ON CONFLICT (user_id, author_id) DO NOTHING`,
        userID, authorID)
    if err != nil {
        return fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
    }
    return nil
}

func (ur *userRepository) GetSubscriptionsByUserID(ctx context.Context, userID int) ([]int, error) {
    var authorIDs []int
    err := ur.db.SelectContext(ctx, &authorIDs,
        `SELECT author_id FROM subscriptions WHERE user_id = $1`, userID)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
    }
    return authorIDs, nil
}

func (ur *userRepository) DeleteSubscription(ctx context.Context, userID, authorID int) error {
    result, err := ur.db.ExecContext(ctx,
        `DELETE FROM subscriptions WHERE user_id = $1 AND author_id = $2`, userID, authorID)
    if err != nil {
        return fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
    }
    if rows, _ := result.RowsAffected(); rows == 0 {
        return apperrors.ErrSqlNoFoundRows
    }
    return nil
}