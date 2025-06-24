package apperrors


// Ошибки сервиса
var (
	ErrEmptyTitle   = &ServiceError{Message: "title cannot be empty"}
	ErrEmptyContent = &ServiceError{Message: "content cannot be empty"}
	ErrInvalidID    = &ServiceError{Message: "invalid post ID"}
	ErrNotFoundPost = &ServiceError{Message: "post not found"}
	ErrDataBase     = &ServiceError{Message: "database error"}

	ErrSqlNoFoundRows 		 = &RepositoryError{Message: "rows not found"}
	ErrSqlForignKeyViolation = &RepositoryError{Message: "foreign key violation"}
	ErrSqlDataBase 		 	 = &RepositoryError{Message: "database error"}

	ErrConfigEmptyDB_CONN_STR = &CustomError{Message: "DB_CONN_STR is required"}
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

// CustomError представляет ошибки для остальных случаев
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}