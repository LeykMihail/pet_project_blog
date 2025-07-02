package app

import (
    "pet_project_blog/internal/repository"
    "pet_project_blog/internal/services"
    "pet_project_blog/internal/handlers"
    "pet_project_blog/internal/config"
	
    "github.com/jmoiron/sqlx"
    "go.uber.org/zap"
)

type AppLayers struct {
    UserRepo    repository.UserRepository
    PostRepo    repository.PostRepository
    UserService services.UserService
    PostService services.PostService
    PostHandler *handlers.PostHandler
    UserHandler *handlers.UserHandler
}

func InitAppLayers(db *sqlx.DB, logger *zap.Logger, cfg *config.Config) *AppLayers {
    userRepo := repository.NewUserRepository(db)
    postRepo := repository.NewPostRepository(db)
    userService := services.NewUserService(userRepo, logger)
    postService := services.NewPostService(postRepo, logger)
    postHandler := handlers.NewPostHandler(postService, logger)
    userHandler := handlers.NewUserHandler(userService, logger, cfg)

    return &AppLayers{
        UserRepo:    userRepo,
        PostRepo:    postRepo,
        UserService: userService,
        PostService: postService,
        PostHandler: postHandler,
        UserHandler: userHandler,
    }
}