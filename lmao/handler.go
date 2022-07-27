package lmao

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/sirupsen/logrus"
	"github.com/tjhu/lmao/lmao/interaction/command"
	"github.com/tjhu/lmao/lmao/interaction/ping"
)

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

func handleInteraction(event discord.InteractionEvent) (*api.InteractionResponse, error) {
	interaction_type := event.Data.InteractionType()
	switch event.Data.InteractionType() {
	case discord.PingInteractionType:
		return ping.HandlePingInteraction()

	case discord.CommandInteractionType:
		cmd := event.Data.(*discord.CommandInteraction)
		return command.HandleCommandInteraction(cmd)

	default:
		error_message := fmt.Sprint("Unrecognized interaction type: ", interaction_type)
		logrus.Warn(error_message)
		return nil, errors.New(error_message)
	}
}

func HandleInteraction(event discord.InteractionEvent) (*api.InteractionResponse, *HandlerError) {
	logrus.Info("Received interaction event: ", event)
	resp, err := handleInteraction(event)
	if err != nil {
		return resp, NewHandlerError(http.StatusBadRequest, err.Error())
	}
	return resp, nil
}
