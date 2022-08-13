package lmao

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/sirupsen/logrus"
	"github.com/uav-gaming/lmao/lmao/interaction/command"
	"github.com/uav-gaming/lmao/lmao/interaction/ping"
)

// A threadsafe instance of the LMAO discord bot for handling requests.
type LMAO struct {
	client         *api.Client
	public_key     []byte
	application_id discord.AppID
}

// Create a LMAO instance.
// It will check against its build time against the current bot commands registered in discord.
// If the the build time is newer, currently registered commands will be replaced by the new ones.
func NewLMAO(token string, public_key []byte, application_id discord.AppID) *LMAO {
	lmao := LMAO{
		api.NewClient(token),
		public_key,
		application_id,
	}

	// TODO: check for existing commands.
	cmds, err := lmao.client.Commands(application_id)
	if err != nil {
		logrus.Error("Failed to get commands: ", cmds)
		return nil
	}
	logrus.Infof("Existing commands: %+v", cmds)

	// TODO: Register commands.

	return &lmao
}

// Handles a discord interaction event and returns an interaction response.
// It always sends back a discord message response to let the user know what happened.
func (bot *LMAO) HandleInteraction(event discord.InteractionEvent) *api.InteractionResponse {
	logrus.Info("Received interaction event: ", event)
	resp, err := bot.handleInteraction(event)
	if err != nil {
		resp = &api.InteractionResponse{
			Type: api.MessageInteractionWithSource,
			Data: &api.InteractionResponseData{
				Content: option.NewNullableString(err.Error()),
			},
		}
	}
	return resp
}

func (bot *LMAO) handleInteraction(event discord.InteractionEvent) (*api.InteractionResponse, error) {
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

// Return true iff the request compiles with the discord authorization protocol.
// https://discord.com/developers/docs/interactions/receiving-and-responding#security-and-authorization
func (bot *LMAO) VerifyRequest(request Request) bool {
	// TODO: get rid of the case sensitivity.
	signature := request.Headers["x-signature-ed25519"]
	if len(signature) <= 0 {
		logrus.Warn("Http header x-signature-ed25519 not set")
		return false
	}
	decoded_signature, err := hex.DecodeString(signature)
	if err != nil {
		logrus.Warn("Invliad request signature: ", signature)
		return false
	}

	timestamp := request.Headers["x-signature-timestamp"]
	if len(timestamp) <= 0 {
		logrus.Warn("Http header x-signature-timestamp not set")
		return false
	}
	body := request.Body

	return ed25519.Verify(bot.public_key, []byte(timestamp+body), decoded_signature)
}
