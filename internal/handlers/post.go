package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"pet_project_blog/internal/apperrors"
	"pet_project_blog/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PostHandler обрабатывает HTTP запросы для постов
type PostHandler struct {
	postService services.PostService
	logger      *zap.Logger
}

// NewPostHandler создает новый экземпляр PostHandler
func NewPostHandler(postService services.PostService, logger *zap.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		logger:      logger,
	}
}

// getHome обрабатывает GET / (главная страница).
func (h *PostHandler) getHome(c *gin.Context) {
	message := "Добро пожаловать в Blog API! 🚀\n\n" +
		"Доступные эндпоинты:\n" +
		"• GET /posts - Просмотр всех постов\n" +
		"• POST /posts - Создать новый пост\n" +
		"• GET /posts/:id - Получить конкретный пост\n" +
		"• POST /posts/:id/comments - Добавить комментарий к посту\n" +
		"• GET /posts/:id/comments - Получить все комментарии к посту\n\n" +
		"Параметры запроса:\n" +
		"• Используйте ?fields=id,title для фильтрации полей ответа\n" +
		"• Пример: /posts?fields=id,title,created_at\n\n" +
		"Для создания и комментирования постов требуется авторизация через /register, а потом /login .\n\n" +
		"Приятного блогинга! ✨"

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})

	// c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(message))
}

// createPost обрабатывает POST /posts (создание поста).
func (h *PostHandler) createPost(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
    if !exists {
        h.logger.Warn("User ID not found in context", zap.String("handler", "createPost"))
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid title or content format", zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Создаем пост по полученным данным
	post, err := h.postService.CreatePost(ctx, input.Title, input.Content, userID.(int))
	if err != nil {
		if err == apperrors.ErrEmptyTitle {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title must not be empty"})
			return
		}
		if err == apperrors.ErrLengthTitle {
			c.JSON(http.StatusBadRequest, gin.H{"error": "maximum length title exceeded"})
			return
		}
		if err == apperrors.ErrEmptyContent {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content must not be empty"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"post": post,
	})
}

// getPost обрабатывает GET /posts/:id (просмотр поста).
func (h *PostHandler) getPost(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")

	// Преобразование строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid ID format", zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Получаем пост из сервиса по указанному ID
	post, err := h.postService.GetPost(ctx, id)
	if err != nil {
		if err == apperrors.ErrNotFoundPost {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if err == apperrors.ErrInvalidID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// getAllPosts обрабатывает GET /posts (получение всех постов).
func (h *PostHandler) getAllPosts(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем параметры из query string
	fields := c.Query("fields") // например: "id,title"

	posts, err := h.postService.GetAllPosts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	// Фильтруем поля если нужно
	if fields != "" {
		fieldList := strings.Split(fields, ",")
		postsWithFields := make([]gin.H, len(posts))
		for i, post := range posts {
			postsWithFields[i] = gin.H{}
			for _, f := range fieldList {
				switch f {
				case "id":
					postsWithFields[i]["id"] = post.ID
				case "title":
					postsWithFields[i]["title"] = post.Title
				case "content":
					postsWithFields[i]["content"] = post.Content
				case "created_at":
					postsWithFields[i]["created_at"] = post.CreatedAt
				case "user_id":
					postsWithFields[i]["user_id"] = post.UserID
				default:
					h.logger.Warn("Invalid filter to response fields", zap.String("fields", fields))
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter to response fields"})
					return
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"posts": postsWithFields})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// createComments обрабатывает POST /posts/:id/comments (добавление нового комментария для поста).
func (h *PostHandler) createComment(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
    if !exists {
        h.logger.Warn("User ID not found in context", zap.String("handler", "createPost"))
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
	
	var input struct {
		Content string `json:"content" binding:"required"`
	}
	idStr := c.Param("id")

	// Преобразование строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid ID format", zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Получаем данные
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid content format", zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Создаем комментарий по полученным данным
	comment, err := h.postService.CreateComment(ctx, id, input.Content, userID.(int))
	if err != nil {
		if err == apperrors.ErrNotFoundPost {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if err == apperrors.ErrInvalidID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}
		if err == apperrors.ErrEmptyContent {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Comment content cannot be empty"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"comment": comment,
	})
}

// getComments обрабатывает GET /posts/:id/comments (получение всех комментариев поста).
func (h *PostHandler) getComments(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")

	// Преобразование строки в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid ID format", zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	comments, err := h.postService.GetCommentsByPostID(ctx, id)
	if err != nil {
		if err == apperrors.ErrNotFoundPost {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if err == apperrors.ErrInvalidID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
