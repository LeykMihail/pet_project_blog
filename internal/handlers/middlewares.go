package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// AuthMiddleware проверяет наличие и валидность user_id cookie для аутентификации пользователя
func AuthMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Получаем cookie user_id
        userID, err := c.Cookie("user_id")
        if err != nil {
            logger.Warn("No user_id cookie found", zap.Error(err))
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // Преобразуем user_id в int
        id, err := strconv.Atoi(userID)
        if err != nil {
            logger.Warn("Invalid user_id cookie value", zap.String("user_id", userID), zap.Error(err))
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
            c.Abort()
            return
        }

        // Сохраняем user_id в контекст для дальнейшего использования
        c.Set("user_id", id)
        logger.Info("User authenticated successfully", zap.Int("user_id", id))
        c.Next()
    }
}