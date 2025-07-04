package services

import (
	"context"
	"time"

	"pet_project_blog/internal/apperrors"
	"pet_project_blog/internal/config"
	"pet_project_blog/internal/models"
	"pet_project_blog/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService определяет интерфейс для работы с пользователями
type UserService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string, cfg *config.Config) (*models.User, string, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

// userService реализует интерфейс UserService
type userService struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Register регистрирует нового пользователя, хеширует пароль и сохраняет пользователя в БД
func (us *userService) Register(ctx context.Context, email, password string) (*models.User, error) {
	us.logger.Info("Start register new user", zap.String("email", email))

	// Валидация данных
	if err := validUserPassword(us.logger, password); err != nil {
		return nil, err
	}

	// Генерация хеша пароля с помощью bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		us.logger.Error("Error generating password hash", zap.Error(err))
		return nil, err
	}

	user := models.User{
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	// Сохранение пользователя в базу данных через репозиторий
	err = us.userRepo.CreateUser(ctx, &user)
	if err != nil {
		switch err {
		case apperrors.ErrSqlUniqueViolation:
			us.logger.Warn("user with this email already exists", zap.Error(err), zap.String("email", email))
			return nil, err
		default:
			us.logger.Error("Failed to save user to database", zap.Error(err))
			return nil, apperrors.ErrDataBase
		}
	}

	us.logger.Info("Register user successfully", zap.Int("id", user.ID), zap.String("email", email))
	return &user, nil
}

// Login выполняет аутентификацию пользователя по email и паролю
func (us *userService) Login(ctx context.Context, email, password string, cfg *config.Config) (*models.User, string, error) {
	us.logger.Info("Start login user", zap.String("email", email))

	// Валидация пароля пользователя
	if err := validUserPassword(us.logger, password); err != nil {
		return nil, "", err
	}

	// Получение пользователя из базы данных по email
	user, err := us.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		switch err {
		case apperrors.ErrSqlNoFoundRows:
			us.logger.Warn("User not found in database", zap.String("email", email))
			return nil, "", apperrors.ErrNotFoundUser
		default:
			us.logger.Error("Failed to fetch user from database", zap.Error(err))
			return nil, "", apperrors.ErrDataBase
		}
	}

	// Проверка правильности пароля с помощью bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		us.logger.Warn("Invalid user password during login", zap.String("email", email))
		return nil, "", apperrors.ErrInvalidPassword
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &models.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Создаем токен с указанными claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Генерируем токен с помощью секретного ключа
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		us.logger.Error("Failed to generate JWT", zap.Error(err))
		return nil, "", apperrors.ErrJWT
	}
	us.logger.Info("Login user successfully", zap.Int("id", user.ID), zap.String("email", email))
	return user, tokenString, nil
}

// Получить пользователя по ID
func (us *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	us.logger.Info("Start get user by ID", zap.Int("id", id))

	// Валидация ID
	if err := validateID(us.logger, id); err != nil {
		return nil, err
	}

	user, err := us.userRepo.GetUserByID(ctx, id)
	if err != nil {
		switch err {
		case apperrors.ErrSqlNoFoundRows:
			us.logger.Warn("User not found in database", zap.Int("id", id))
			return nil, apperrors.ErrNotFoundUser
		default:
			us.logger.Error("Failed to fetch user from database", zap.Error(err))
			return nil, apperrors.ErrDataBase
		}
	}

	us.logger.Info("Get user by ID successfully", zap.Int("id", user.ID), zap.String("email", user.Email))
	return user, nil
}
