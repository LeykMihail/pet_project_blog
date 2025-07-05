package handlers

import (
	"net/http"

	"pet_project_blog/internal/apperrors"
	"pet_project_blog/internal/config"
	"pet_project_blog/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler обрабатывает HTTP запросы для пользователей
type UserHandler struct {
	userService services.UserService
	logger      *zap.Logger
	cfg         *config.Config
}

// NewPostHandler создает новый экземпляр PostHandler
func NewUserHandler(userService services.UserService, logger *zap.Logger, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
		cfg:         cfg,
	}
}

// register обрабатывает POST /register (регистрация пользователя).
func (h *UserHandler) register(c *gin.Context) {
	ctx := c.Request.Context()
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid input format", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.userService.Register(ctx, input.Email, input.Password)
	if err != nil {
		switch err {
		case apperrors.ErrEmptyPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must not be empty"})
			return
		case apperrors.ErrLenghtPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": "password length must be between 8 and 64 characters"})
			return
		case apperrors.ErrSqlUniqueViolation:
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"user": gin.H{"id": user.ID, "email": user.Email}})
}

// login обрабатывает POST /login (авторизация пользователя).
func (h *UserHandler) login(c *gin.Context) {
	ctx := c.Request.Context()
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid input format", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, tokenString, err := h.userService.Login(ctx, input.Email, input.Password, h.cfg)
	if err != nil {
		switch err {
		case apperrors.ErrEmptyPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "password must not be empty"})
			return
		case apperrors.ErrLenghtPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "password length must be between 8 and 64 characters"})
			return
		case apperrors.ErrNotFoundUser:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		case apperrors.ErrInvalidPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		case apperrors.ErrDataBase:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authorized user"})
			return
		}
	}

	c.Header("Authorization", "Bearer "+tokenString)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
		"user_id": user.ID,
		"token":   tokenString,
	})
}
