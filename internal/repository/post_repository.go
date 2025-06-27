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

// PostRepository определяет интерфейс для работы с постами в базе данных
type PostRepository interface {
	CreatePost(ctx context.Context, post *models.Post) error
	GetPost(ctx context.Context, id int) (*models.Post, error)
	GetAllPosts(ctx context.Context) ([]*models.Post, error)

	CreatComment(ctx context.Context, comment *models.Comment) error
	GetCommentsByPostID(ctx context.Context, id int) ([]*models.Comment, error)
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
		`INSERT INTO posts (title, content, created_at, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		post.Title, post.Content, post.CreatedAt, post.UserID,
	).Scan(&id)

	if err != nil {
		if fkErr := checkErrForeignKeyViolation(err); fkErr != nil {
			return fkErr
		}
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
		`SELECT id, title, content, created_at, user_id FROM posts WHERE id = $1`, id)
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
		`SELECT id, title, content, created_at, user_id FROM posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}

	return posts, nil
}

func (pr *postRepository) CreatComment(ctx context.Context, comment *models.Comment) error {
	// Выполнение SQL запроса для вставки нового поста с возвратом ID
	var commentID int
	err := pr.db.QueryRowContext(ctx,
		`INSERT INTO comments (post_id, content, created_at, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		comment.PostID, comment.Content, comment.CreatedAt, comment.UserID,
	).Scan(&commentID)

	if err != nil {
		if fkErr := checkErrForeignKeyViolation(err); fkErr != nil {
			return fkErr
		}
		return fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}

	// Устанавливаем ID созданного поста
	comment.ID = commentID
	return nil
}

func (pr *postRepository) GetCommentsByPostID(ctx context.Context, id int) ([]*models.Comment, error) {
	var comments []*models.Comment

	// Выполнение SQL запроса для получения всех комментариев поста
	err := pr.db.SelectContext(ctx, &comments,
		`SELECT id, post_id, content, created_at, user_id FROM comments WHERE post_id = $1 ORDER BY created_at DESC`, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperrors.ErrSqlDataBase, err)
	}

	return comments, nil
}
