package handlers

import (
	"pet_project_blog/internal/config"
	"pet_project_blog/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes регистрирует маршруты для постов.
func RegisterRoutes(r *gin.Engine, cfg *config.Config, postHandler *PostHandler, userHandler *UserHandler, userService services.UserService) {
	r.GET("/", postHandler.getHome)
	r.GET("/posts", postHandler.getAllPosts)
	r.POST("/posts", AuthMiddleware(postHandler.logger, cfg, userService), postHandler.createPost)
	r.GET("/posts/:id", postHandler.getPost)
	r.POST("/posts/:id/comments", AuthMiddleware(postHandler.logger, cfg, userService), postHandler.createComment)
	r.GET("/posts/:id/comments", postHandler.getComments)

	r.POST("/register", userHandler.register)
	r.POST("/login", userHandler.login)
}
