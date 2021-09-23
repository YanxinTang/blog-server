package e

type SessionError struct {
	Message string
}

func (se SessionError) Error() string {
	return se.Message
}

func NewSessionError(code int, msg string) SessionError {
	return SessionError{msg}
}
