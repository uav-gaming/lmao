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

	return nil, NewHandlerError(http.StatusOK, fmt.Sprint("Unrecognized interaction type: ", interaction_type))
}
