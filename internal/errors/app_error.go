package errors

type AppError struct {
	Message			string
	StatusCode		int
	Data			string
}

func (e *AppError) Error() string {
	return e.Message	
}

func BadRequest(message string) *AppError {
	return &AppError{
		Message: 	message,
		StatusCode: 400,
		Data:		"BAD_REQUEST",
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Message: 	message,
		StatusCode: 404,
		Data:		"NOT_FOUND",
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Message: 	message,
		StatusCode: 401,
		Data: 		"UNAUTHORIZED",
	}
}

func InternalServerError(message string) *AppError {
	return &AppError{
		Message: 	message,
		StatusCode: 500,
		Data:		"INTERNAL_SERVER_ERROR",
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Message: message,
		StatusCode: 409,
		Data: "CONFLICT",
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Message: message,
		StatusCode: 403,
		Data: "FORBIDDEN",
	}
}