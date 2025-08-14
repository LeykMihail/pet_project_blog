package apperrors


// Ошибки сервиса
var (
	ErrEmptyTitle      = &ServiceError{Message: "title cannot be empty"}
	ErrLengthTitle	   = &ServiceError{Message: "maximum length title exceeded"}
	ErrEmptyContent    = &ServiceError{Message: "content cannot be empty"}
	ErrEmptyPassword   = &ServiceError{Message: "password cannot be empty"}
	ErrLenghtPassword  = &ServiceError{Message: "incorrect password length"}
	ErrInvalidID       = &ServiceError{Message: "invalid ID"}
	ErrInvalidPassword = &ServiceError{Message: "invalid user password"}
	ErrNotFoundPost    = &ServiceError{Message: "post not found"}
	ErrNotFoundComment = &ServiceError{Message: "comment not found"}
	ErrNotFoundUser    = &ServiceError{Message: "user not found"}
	ErrDataBase        = &ServiceError{Message: "database error"}
	ErrJWT			   = &ServiceError{Message: "JWT error"}

	ErrSqlNoFoundRows 		 = &RepositoryError{Message: "rows not found"}
	ErrSqlForignKeyViolation = &RepositoryError{Message: "foreign key violation"}
	ErrSqlUniqueViolation    = &RepositoryError{Message: "unique violation"} 
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