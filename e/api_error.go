package e

type ApiError interface {
	Code() int
	Error() string
}

type apiError struct {
	code    int
	message string
}

func (e apiError) Error() string {
	return e.message
}

func (e apiError) Code() int {
	return e.code
}

func New(code int, msg string) ApiError {
	return apiError{code, msg}
}
