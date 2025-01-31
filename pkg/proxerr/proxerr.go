package proxerr

type Error interface {
	Error() string
	Unwrap() error
}

type proxyError struct {
	background error
	message    string
}

func (e *proxyError) Error() string {
	return e.message
}

func (e *proxyError) Unwrap() error {
	return e.background
}

func New(background error, message string) Error {
	return &proxyError{
		background: background,
		message:    message,
	}
}
