package lmao

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/sirupsen/logrus"
	"github.com/uav-gaming/lmao/lmao/commands"
)

// A threadsafe instance of the LMAO discord bot for handling requests.
type LMAO struct {
	client         *api.Client
	public_key     ed25519.PublicKey
	application_id discord.AppID
	cmds           commands.CommandRegistrar
}

// Create a LMAO instance.
// It will check against its build time against the current bot commands registered in discord.
// If the the build time is newer, currently registered commands will be replaced by the new ones.
func NewLMAO(token string, public_key ed25519.PublicKey, application_id discord.AppID) (*LMAO, error) {
	logrus.Info("Initilizing lmao with public key: ", hex.EncodeToString(public_key), " and app_id: ", application_id)
	lmao := LMAO{
		api.NewClient("Bot " + token),
		public_key,
		application_id,
		commands.DefaultCommandRegistrar(),
	}

	// Check for existing commands.
	// And register/update commands if necessary.
	// Updating means deleting existing commands and registering new ones for simplicity.
	var cmdsToDelete []*discord.Command
	cmds, err := lmao.client.Commands(application_id)
	if err != nil {
		logrus.Warn("Failed to get commands: ", err, ". Attempting to register for new ones.")
	} else {
		logrus.Infof("Existing commands: %+v", cmds)
		// TODO: check if commands are acutally needed to be updated.
		for _, cmd := range cmds {
			cmdsToDelete = append(cmdsToDelete, &cmd)
		}
	}

	// Delete commands.
	for _, cmd := range cmdsToDelete {
		logrus.Infof("Deleting command: %+v", *cmd)
		err := lmao.client.DeleteCommand(application_id, cmd.ID)
		if err != nil {
			return nil, errors.New(fmt.Sprint("failed to delete command ", cmd.ID, ": ", err.Error()))
		}
	}

	// Register commands
	for _, cmd := range lmao.cmds.Commands {
		logrus.Infof("Registering commands: %+v", cmd)

		_, err := lmao.client.CreateCommand(application_id, cmd)
		if err != nil {
			return nil, errors.New("failed to create command: " + err.Error())
		}
	}

	return &lmao, nil
}

// Handles a discord interaction event and returns an interaction response.
// It always sends back a discord message response to let the user know what happened.
func (bot *LMAO) HandleInteraction(event discord.InteractionEvent) *api.InteractionResponse {
	logrus.WithField("interaction", event).Info("Handling interaction event")
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
	if event.AppID != bot.application_id {
		logrus.Error("Incorrect application id: ", event.AppID, " vs ", bot.application_id)
	}
	interaction_type := event.Data.InteractionType()
	switch event.Data.InteractionType() {
	case discord.PingInteractionType:
		logrus.Info("Received ping. Responding with pong.")
		return &api.InteractionResponse{
			Type: api.PongInteraction,
		}, nil

	case discord.CommandInteractionType:
		cmd := event.Data.(*discord.CommandInteraction)
		return bot.cmds.HandleCommand(cmd)

	default:
		error_message := fmt.Sprint("Unrecognized interaction type: ", interaction_type)
		logrus.Warn(error_message)
		return nil, errors.New(error_message)
	}
}

// Return true iff the request compiles with the discord authorization protocol.
// https://discord.com/developers/docs/interactions/receiving-and-responding#security-and-authorization
func (bot *LMAO) VerifyRequest(request Request) bool {
	// Prepare signature.
	// TODO: get rid of the case sensitivity.
	signature := request.Headers["x-signature-ed25519"]
	if len(signature) <= 0 {
		logrus.Warn("Http header x-signature-ed25519 not set")
		return false
	}
	sig, err := hex.DecodeString(signature)
	if err != nil {
		logrus.Warn("Non-hex request signature: ", signature)
		return false
	}
	if len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
		logrus.Warn("Invalid ed25519 signature format: ", signature)
		return false
	}

	// Prepare message.
	var msg bytes.Buffer
	timestamp := request.Headers["x-signature-timestamp"]
	if len(timestamp) <= 0 {
		logrus.Warn("Http header x-signature-timestamp not set")
		return false
	}
	msg.WriteString(timestamp)
	msg.WriteString(request.Body)

	// Verify the signature.
	return ed25519.Verify(bot.public_key, msg.Bytes(), sig)
}
