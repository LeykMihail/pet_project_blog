package repository

import (
	"context"
	"database/sql"
	"fmt"

	"pet_project_blog/internal/models"
	"pet_project_blog/internal/apperrors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// PostRepository определяет интерфейс для работы с постами в базе данных
type PostRepository interface {
	CreatePost(ctx context.Context, post *models.Post) error
	GetPost(ctx context.Context, id int) (*models.Post, error)
	GetAllPosts(ctx context.Context) ([]*models.Post, error)
	// другие методы
}

// postRepository реализует интерфейс PostRepository
type postRepository struct {
	db *sqlx.DB
}

// NewPostRepository создает новый экземпляр PostRepository
func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

// CreatePost сохраняет новый пост в базу данных
func (pr *postRepository) CreatePost(ctx context.Context, post *models.Post) error {
	// Выполнение SQL запроса для вставки нового поста с возвратом ID
	var id int
	err := pr.db.QueryRowContext(ctx,
		`INSERT INTO posts (title, content, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`,
		post.Title, post.Content, post.CreatedAt,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}

	// Устанавливаем ID созданного поста
	post.ID = id
	return nil
}

// GetPost извлекает пост из базы данных по его ID
func (pr *postRepository) GetPost(ctx context.Context, id int) (*models.Post, error) {
	var post models.Post
	// Выполняем SQL запрос для получения поста по ID
	err := pr.db.GetContext(ctx, &post,
		`SELECT id, title, content, created_at FROM posts WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrSqlNoFoundRows
		}
		return nil, fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}
	return &post, nil
}

// GetAllPosts извлекает все посты из базы данных
func (pr *postRepository) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	var posts []*models.Post
	// Выполнение SQL запроса для получения всех постов
	err := pr.db.SelectContext(ctx, &posts,
		`SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}

	return posts, nil
}
