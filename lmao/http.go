package lmao

import (
	"encoding/json"
	"net/http"

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
		StatusCode: http.StatusOK,
		Body:       string(buffer),
	}, nil
}
