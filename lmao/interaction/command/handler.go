package command

import (
	"errors"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/sirupsen/logrus"
)

func HandleCommandInteraction(command *discord.CommandInteraction) (*api.InteractionResponse, error) {
	logrus.Info("Received command: ", command)
	return nil, errors.New("can't handle commands atm")
}
