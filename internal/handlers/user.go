package handlers

import (
	"net/http"
	"strconv"

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

// createSubscription обрабатывает POST /subscriptions (подписка на автора).
func (h *UserHandler) createSubscription(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		// Пользователь не авторизован
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authorIDStr := c.Query("author_id")
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		// Некорректный формат ID автора
		h.logger.Warn("Invalid author ID format", zap.String("author_id", authorIDStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}
	
	// Пытаемся создать подписку
	err = h.userService.CreateSubscription(ctx, userID.(int), authorID)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidID:
			// Некорректный ID пользователя или автора
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user or author ID"})
			return
		case apperrors.ErrNotFoundUser:
			// Автор не найден
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
			return
		case apperrors.ErrSelfSubscription:
			// Попытка подписаться на самого себя
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot subscribe to yourself"})
			return
		default:
			// Другая ошибка при создании подписки
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
			return
		}
	}
	// Подписка успешно создана
	c.JSON(http.StatusCreated, gin.H{"message": "Subscription created"})
}

// getSubscriptions обрабатывает GET /subscriptions (получение списка подписок пользователя).
func (h *UserHandler) getSubscriptions(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		// Пользователь не авторизован
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authorIDs, err := h.userService.GetSubscriptionsByUserID(ctx, userID.(int))
	if err != nil {
		// Ошибка при получении подписок пользователя
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
		return
	}
	// Успешно возвращаем список подписок
	c.JSON(http.StatusOK, gin.H{"subscriptions": authorIDs})
}

// deleteSubscription обрабатывает DELETE /subscriptions/:authorID (удаление подписки пользователя на автора).
func (h *UserHandler) deleteSubscription(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		// Пользователь не авторизован
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authorIDStr := c.Param("authorId")
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		// Некорректный формат authorId в параметре запроса
		h.logger.Warn("Invalid author ID format", zap.String("author_id", authorIDStr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}
	// Пытаемся удалить подписку пользователя на автора
	err = h.userService.DeleteSubscription(ctx, userID.(int), authorID)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidID:
			// Некорректный ID пользователя или автора
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user or author ID"})
			return
		case apperrors.ErrNotFoundSubscription:
			// Подписка не найдена
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		default:
			// Другая ошибка при удалении подписки
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
			return
		}
	}
	// Подписка успешно удалена
	c.JSON(http.StatusNoContent, gin.H{})
}
