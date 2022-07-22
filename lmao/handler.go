package lmao

import (
	"fmt"
	"net/http"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/sirupsen/logrus"
)

func HandleEvent(event discord.InteractionEvent) (*api.InteractionResponse, *HandlerError) {
	interaction_type := event.Data.InteractionType()
	if event.Data.InteractionType() == discord.PingInteractionType {
		logrus.Info("Received ping. Responding with pong.")
		return &api.InteractionResponse{
			Type: api.PongInteraction,
		}, nil
	}

	error_message := fmt.Sprint("Unrecognized interaction type: ", interaction_type)
	logrus.Warn(error_message)
	return nil, NewHandlerError(http.StatusBadRequest, error_message)
}
