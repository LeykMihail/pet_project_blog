package services

import (
	"context"
	"time"
	
	"pet_project_blog/internal/models"
	"pet_project_blog/internal/repository"
	"pet_project_blog/internal/apperrors"

	"go.uber.org/zap"
)

// PostService определяет интерфейс для бизнес-логики работы с постами
type PostService interface {
	CreatePost(ctx context.Context, title, content string) (*models.Post, error)
	GetPost(ctx context.Context, id int) (*models.Post, error)
	GetAllPosts(ctx context.Context) ([]*models.Post, error)
}

// postService реализует интерфейс PostService
type postService struct {
	postRepo repository.PostRepository
	logger   *zap.Logger
}

// NewPostService создает новый экземпляр PostService
func NewPostService(postRepo repository.PostRepository, logger *zap.Logger) PostService {
	return &postService{
		postRepo: postRepo, logger: logger,
	}
}

// CreatePost создает новый пост с валидацией и бизнес-логикой
func (ps *postService) CreatePost(ctx context.Context, title, content string) (*models.Post, error) {
	ps.logger.Info("Creating new post", zap.String("title", title))

	// Валидация входных данных
	if title == "" {
		ps.logger.Warn("Empty title when creating post")
		return nil, apperrors.ErrEmptyTitle
	}
	if content == "" {
		ps.logger.Warn("Empty content when creating post")
		return nil, apperrors.ErrEmptyContent
	}

	// Создание модели поста с временной меткой
	post := models.Post{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	// Сохранение в базу данных через репозиторий
	err := ps.postRepo.CreatePost(ctx, &post)
	if err != nil {
		ps.logger.Error("Failed to save post to database", zap.Error(err))
		return nil, apperrors.ErrDataBase
	}

	ps.logger.Info("Post created successfully", zap.Int("id", post.ID))
	return &post, nil
}

// GetPost получает пост по ID с обработкой ошибок
func (ps *postService) GetPost(ctx context.Context, id int) (*models.Post, error) {
	ps.logger.Info("Fetching post by ID", zap.Int("id", id))

	// Валидация ID
	if id <= 0 {
		ps.logger.Warn("Invalid ID", zap.Int("id", id))
		return nil, apperrors.ErrInvalidID
	}

	// Получение поста из базы данных через репозиторий
	post, err := ps.postRepo.GetPost(ctx, id)
	if err != nil {
		if err == apperrors.ErrSqlNoFoundRows {
			ps.logger.Warn("Post not found in database", zap.Int("id", id))
			return nil, apperrors.ErrNotFoundPost
		}
		ps.logger.Error("Failed to fetch post from database", zap.Error(err))
        return nil, apperrors.ErrDataBase
	}

	ps.logger.Info("Fetching post successfully", zap.Int("id", post.ID))
	return post, nil
}

// GetAllPosts получает все посты с обработкой ошибок
func (ps *postService) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	ps.logger.Info("Fetching all posts")

	// Получение всех постов из базы данных через репозиторий
	posts, err := ps.postRepo.GetAllPosts(ctx)
	if err != nil {
		ps.logger.Error("Failed to fetch all posts from database", zap.Error(err))
		return nil, apperrors.ErrDataBase
	}

	ps.logger.Info("Fetched all posts successfully", zap.Int("count", len(posts)))
	return posts, nil
}
