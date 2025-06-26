package services

import (
	"pet_project_blog/internal/apperrors"
	
	"go.uber.org/zap"
)

// Функция валидции id
func validateID(logger *zap.Logger, id int) error {
	if id <= 0 {
		logger.Warn("Invalid ID", zap.Int("id", id))
		return apperrors.ErrInvalidID
	}
	return nil
}

// Функция валидации title
func validatePostTitle(logger *zap.Logger, title string) error {
	if title == "" {
		logger.Warn("Empty title when creating post")
		return apperrors.ErrEmptyTitle
	}
	if len(title) > 100 {
		logger.Warn("Maximum length title exceeded")
		return apperrors.ErrLengthTitle
	}
	return nil
}

// Функция валидации content у post
func validatePostContent(logger *zap.Logger, content string) error {
	if content == "" {
		logger.Warn("Empty content when creating post")
		return apperrors.ErrEmptyContent
	}
	return nil
}

// Функция валидации content у comment
func validateCommentContent(logger *zap.Logger, content string) error {
	if content == "" {
		logger.Warn("Empty content when creating comment")
		return apperrors.ErrEmptyContent
	}
	return nil
}
