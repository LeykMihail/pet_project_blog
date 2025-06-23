package apperrors


// Ошибки сервиса
var (

	ErrEmptyTitle   = &ServiceError{Message: "title cannot be empty"}
	ErrEmptyContent = &ServiceError{Message: "content cannot be empty"}
	ErrInvalidID    = &ServiceError{Message: "invalid post ID"}
	ErrNotFoundPost = &ServiceError{Message: "post not found"}
	ErrDataBase     = &ServiceError{Message: "database error"}

	ErrSqlNoFoundRows 	= &RepositoryError{Message: "rows not found in database"}
	ErrSqlDataBase 		= &RepositoryError{Message: "database error"}
)

// ServiceError представляет ошибку сервиса
type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// RepositoryError представляет ошибку репозитория
type RepositoryError struct {
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}
