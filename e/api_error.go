package e

type ApiError struct {
	Code    int
	Message string
}

func (ae ApiError) Error() string {
	return ae.Message
}

func New(code int, msg string) ApiError {
	return ApiError{code, msg}
}
