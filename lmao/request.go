package lmao

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/diamondburned/arikawa/utils/json"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Interaction struct {
	Type discord.InteractionDataType `json:"type"`
	Data json.Raw                    `json:"data"`
}

type Request struct {
	Headers        map[string]interface{}                 `json:"headers"`
	Body           string                                 `json:"body"`
	RequestContext events.LambdaFunctionURLRequestContext `json:"requestContext"`
}
