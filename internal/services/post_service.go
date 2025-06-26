package services

import (
	"context"
	"time"

	"pet_project_blog/internal/apperrors"
	"pet_project_blog/internal/models"
	"pet_project_blog/internal/repository"

	"go.uber.org/zap"
)

// PostService определяет интерфейс для бизнес-логики работы с постами
type PostService interface {
	CreatePost(ctx context.Context, title, content string) (*models.Post, error)
	GetPost(ctx context.Context, id int) (*models.Post, error)
	GetAllPosts(ctx context.Context) ([]*models.Post, error)

	CreateComment(ctx context.Context, postID int, content string) (*models.Comment, error)
	GetCommentsByPostID(ctx context.Context, id int) ([]*models.Comment, error)
}

// postService реализует интерфейс PostService
type postService struct {
	postRepo repository.PostRepository
	logger   *zap.Logger
}

// NewPostService создает новый экземпляр PostService
func NewPostService(postRepo repository.PostRepository, logger *zap.Logger) PostService {
	return &postService{
		postRepo: postRepo, 
		logger: logger,
	}
}

// CreatePost создает новый пост с валидацией и бизнес-логикой
func (ps *postService) CreatePost(ctx context.Context, title, content string) (*models.Post, error) {
	ps.logger.Info("Start creating new post", zap.String("title", title))

	// Валидация входных данных
	if err := validatePostTitle(ps.logger, title); err != nil {

		return nil, err
	}
	if err := validatePostContent(ps.logger, content); err != nil {
		return nil, err
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
	ps.logger.Info("Start fetching post by ID", zap.Int("id", id))

	// Валидация ID
	if err := validateID(ps.logger, id); err != nil {
		return nil, err
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

	// Получение всех комментариев поста из базы данных через репозиторий
	comments, err := ps.postRepo.GetCommentsByPostID(ctx, post.ID)
	if err != nil {
		ps.logger.Error("Failed to fetch all comments from database", zap.Error(err), zap.Int("postID", post.ID))
		return nil, apperrors.ErrDataBase
	}
	post.Comments = comments

	ps.logger.Info("Fetching post successfully", zap.Int("id", post.ID))
	return post, nil
}

// GetAllPosts получает все посты с обработкой ошибок
func (ps *postService) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	ps.logger.Info("Start fetching all posts")

	// Получение всех постов из базы данных через репозиторий
	posts, err := ps.postRepo.GetAllPosts(ctx)
	if err != nil {
		ps.logger.Error("Failed to fetch all posts from database", zap.Error(err))
		return nil, apperrors.ErrDataBase
	}

	ps.logger.Info("Fetched all posts successfully", zap.Int("count", len(posts)))
	return posts, nil
}

func (ps *postService) CreateComment(ctx context.Context, postID int, content string) (*models.Comment, error) {
	ps.logger.Info("Start creating new comment", zap.Int("post ID", postID))

	// Валидация данных
	if err := validateID(ps.logger, postID); err != nil {
		return nil, err
	}
	if err := validateCommentContent(ps.logger, content); err != nil {
		return nil, err
	}

	// Создание модели поста с временной меткой
	comment := models.Comment{
		PostID:    postID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	// Сохранение в базу данных через репозиторий
	err := ps.postRepo.CreatComment(ctx, &comment)
	if err != nil {
		if err == apperrors.ErrSqlForignKeyViolation {
			ps.logger.Warn("Post not found in database", zap.Error(err))
			return nil, apperrors.ErrNotFoundPost
		}
		ps.logger.Error("Failed to save comment to database", zap.Error(err), zap.Int("postID", postID))
		return nil, apperrors.ErrDataBase
	}

	ps.logger.Info("Comment created successfully", zap.Int("id", comment.ID))
	return &comment, nil
}

func (ps *postService) GetCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error) {
	ps.logger.Info("Start fetching all comments", zap.Int("postID", postID))

	// Валидация ID
	if err := validateID(ps.logger, postID); err != nil {
		return nil, err
	}
	
	// Проверка существования поста
	_, err := ps.postRepo.GetPost(ctx, postID)
	if err != nil {
		if err == apperrors.ErrSqlNoFoundRows {
			ps.logger.Warn("Post not found in database", zap.Int("id", postID))
			return nil, apperrors.ErrNotFoundPost
		}
		ps.logger.Error("Failed to fetch post from database", zap.Error(err))
        return nil, apperrors.ErrDataBase
	}

	// Получение всех комментариев поста из базы данных через репозиторий
	comments, err := ps.postRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		ps.logger.Error("Failed to fetch all comments from database", zap.Error(err), zap.Int("postID", postID))
		return nil, apperrors.ErrDataBase
	}

	ps.logger.Info("Fetched all comments successfully", zap.Int("count", len(comments)), zap.Int("postID", postID))
	return comments, nil
}
