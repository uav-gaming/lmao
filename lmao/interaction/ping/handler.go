package ping

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/sirupsen/logrus"
)

func HandlePingInteraction() (*api.InteractionResponse, error) {
	logrus.Info("Received ping. Responding with pong.")
	return &api.InteractionResponse{
		Type: api.PongInteraction,
	}, nil
}
