package lmao

import "fmt"

type LMAOError struct {
	StatusCode int
	Message    string
}

func NewLMAOError(status_code int, message string) *LMAOError {
	return &LMAOError{
		StatusCode: status_code,
		Message:    message,
	}
}

func (e *LMAOError) Error() string {
	return fmt.Sprint(e.StatusCode, ": ", e.Message)
}

func (e *LMAOError) ToResponse() Response {
	return Response{
		StatusCode: e.StatusCode,
		Body:       fmt.Sprintf("{message: %s}", e.Message),
	}
}
