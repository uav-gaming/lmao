package lmao

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/sirupsen/logrus"
	"github.com/tjhu/discord_api"
	"github.com/tjhu/lmao/lmao/interaction/command"
	"github.com/tjhu/lmao/lmao/interaction/ping"
)

// A threadsafe instance of the LMAO discord bot for handling requests.
type LMAO struct {
	discord_public_key []byte
	da                 *discord_api.DiscordApi
}

// Create a LMAO instance.
// It will check against its build time against the current bot commands registered in discord.
// If the the build time is newer, currently registered commands will be replaced by the new ones.
func NewLMAO() *LMAO {
	// TODO: check for existing commands.

	// TODO: Register commands.

	return &LMAO{}
}

// Handles a discord interaction event and returns an interaction response.
func (bot *LMAO) HandleInteraction(event discord.InteractionEvent) (*api.InteractionResponse, *LMAOError) {
	logrus.Info("Received interaction event: ", event)
	resp, err := bot.handleInteraction(event)
	if err != nil {
		return resp, NewLMAOError(http.StatusBadRequest, err.Error())
	}
	return resp, nil
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

	return ed25519.Verify(bot.discord_public_key, []byte(timestamp+body), decoded_signature)
}
