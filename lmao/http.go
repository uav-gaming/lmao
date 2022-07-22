package lmao

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/diamondburned/arikawa/v3/api"
)

// type Request struct {
// 	Headers        map[string]interface{}                 `json:"headers"`
// 	Body           string                                 `json:"body"`
// 	RequestContext events.LambdaFunctionURLRequestContext `json:"requestContext"`
// }

type Request = events.LambdaFunctionURLRequest

type Response = events.LambdaFunctionURLResponse

func ToHttpResponse(r *api.InteractionResponse) (Response, error) {
	buffer, err := json.Marshal(r)
	if err != nil {
		return Response{}, err
	}
	return Response{
		StatusCode: 200,
		Body:       string(buffer),
	}, nil
}

type HandlerError struct {
	StatusCode int
	Message    string
}

func NewHandlerError(status_code int, message string) *HandlerError {
	return &HandlerError{
		StatusCode: status_code,
		Message:    message,
	}
}

func (e *HandlerError) Error() string {
	return fmt.Sprint(e.StatusCode, ": ", e.Message)
}

func (e *HandlerError) ToResponse() Response {
	return Response{
		StatusCode: e.StatusCode,
		Body:       fmt.Sprintf("{message: %s}", e.Message),
	}
}
