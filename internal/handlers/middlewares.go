package handlers

import (
    "net/http"
    "strings"
    "fmt"
    
    "pet_project_blog/internal/config"
    "pet_project_blog/internal/models"
    "pet_project_blog/internal/services"
    "pet_project_blog/internal/apperrors"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "go.uber.org/zap"
)

// AuthMiddleware — middleware для проверки JWT авторизации пользователя
func AuthMiddleware(logger *zap.Logger, cfg *config.Config, us services.UserService) gin.HandlerFunc {
    return func(c *gin.Context) {
        logger.Info("Running Auth Middleware to verify user authorization")
        // Получаем заголовок Authorization из запроса
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            // Если заголовок отсутствует — возвращаем 401
            logger.Warn("No Authorization header found")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        // Извлекаем токен из заголовка (убираем префикс "Bearer ")
        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenStr == "" {
            logger.Warn("No token found in Authorization header")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // Создаём структуру claims для парсинга токена
        claim := models.Claims{}

        // Парсим токен и валидируем claims
        token, err := jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (interface{}, error) {
            // Проверяем, что используется ожидаемый метод подписи (HMAC)
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                logger.Warn("Unexpected signing method", zap.String("method", token.Header["alg"].(string)))
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"].(string))
            }
            // Возвращаем секретный ключ для проверки подписи
            return []byte(cfg.JWTSecret), nil
        })
        if err != nil {
            logger.Warn("Invalid or malformed token", zap.Error(err))
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Проверяем валидность токена и наличие user_id
        if !token.Valid {
            logger.Warn("Token invalid")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Проверяем, что user_id присутствует и валиден
        userID := claim.UserID
        // Проверяем, что пользователь с таким userID существует в базе данных
        _, err = us.GetUserByID(c.Request.Context(), userID)
        if err != nil {
            if err == apperrors.ErrNotFoundUser {
                logger.Warn("User not found for user_id in JWT", zap.Int("user_id", userID), zap.Error(err))
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            }else{
                logger.Error("Database error while checking user existence", zap.Error(err))
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            }
            c.Abort()
            return
        }

        // Добавляем user_id в контекст Gin для дальнейшего использования в хендлерах
        c.Set("user_id", userID)
        logger.Info("The user has successfully authenticated via JWT", zap.Int("user_id", userID))

        // Передаём управление следующему обработчику
        c.Next()
    }
}