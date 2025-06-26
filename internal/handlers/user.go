package handlers

import (
	"net/http"
	"strconv"

	"pet_project_blog/internal/services"
	"pet_project_blog/internal/apperrors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler обрабатывает HTTP запросы для пользователей
type UserHandler struct {
	userService services.UserService
	logger      *zap.Logger
}

// NewPostHandler создает новый экземпляр PostHandler
func NewUserHandler(userService services.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func RegisterRoutesUser(r *gin.Engine, userHandler *UserHandler) {
	r.POST("/register", userHandler.register)
    r.POST("/login", userHandler.login)
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
        if err == apperrors.ErrEmptyPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must not be empty"})
			return
		}
		if err == apperrors.ErrLenghtPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password length must be between 8 and 64 characters"})
			return
		}
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
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

    user, err := h.userService.Login(ctx, input.Email, input.Password)
    if err != nil {
		if err == apperrors.ErrEmptyPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "password must not be empty"})
			return
		}
		if err == apperrors.ErrLenghtPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "password length must be between 8 and 64 characters"})
			return
		}
        if err == apperrors.ErrNotFoundUser {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            return
        }
        if err == apperrors.ErrInvalidPassword {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
            return
        }
        if err == apperrors.ErrDataBase {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            return
        }
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to authorized user"})
        return
    }
    // Устанавливаем cookie с user_id для идентификации пользователя в сессии
    c.SetCookie("user_id", strconv.Itoa(user.ID), 3600, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "Logged in", "user_id": user.ID})
}