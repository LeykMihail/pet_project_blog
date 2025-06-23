package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"pet_project_blog/internal/services"
	"pet_project_blog/internal/apperrors"

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

// RegisterRoutes регистрирует маршруты для постов.
func RegisterRoutes(r *gin.Engine, postHandler *PostHandler) {
	r.GET("/", postHandler.getHome)
	r.GET("/posts", postHandler.getAllPosts)
	r.POST("/posts", postHandler.createPost)
	r.GET("/posts/:id", postHandler.getPost)
}

// getHome обрабатывает GET / (главная страница).
func (h *PostHandler) getHome(c *gin.Context) {
	message := "Welcome to the Blog API! 🚀\n\n" +
		"Available endpoints:\n" +
		"• GET /posts - View all posts\n" +
		"• POST /posts - Create a new post\n" +
		"• GET /posts/:id - Get a specific post\n\n" +
		"Query Parameters:\n" +
		"• Use ?fields=id,title to filter response fields\n" +
		"• Example: /posts?fields=id,title,created_at\n\n" +
		"Happy blogging! ✨"

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})

	// c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(message))
}

// createPost обрабатывает POST /posts (создание поста).
func (h *PostHandler) createPost(c *gin.Context) {
	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid title or content format", zap.String("id", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Создаем пост по полученным данным
	post, err := h.postService.CreatePost(ctx, input.Title, input.Content)
	if err != nil {
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
		if errors.Is(err, apperrors.ErrNotFoundPost) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if errors.Is(err, apperrors.ErrInvalidID) {
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
