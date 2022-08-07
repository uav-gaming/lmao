package discordapi

import (
	"fmt"
	"net/http"

	"github.com/diamondburned/arikawa/v3/discord"
)

const DISCORD_APPLICATION_URL string = "https://discord.com/api/v10/applications"

// Post request to "`application_url`/`subendpoint`" with `body` as the json body.
func (da *DiscordApi) SendApplicationRequest(method string, subendpoint string, body []byte) error {
	url := fmt.Sprintf("%s/%s", da.application_url, subendpoint)
	return da.SendRequest(method, url, body)
}

func (da *DiscordApi) SendApplicationCommandRequest(method string, body []byte) error {
	return da.SendApplicationRequest(method, "commands", body)
}

func (da *DiscordApi) RegsiterCommand(name string, description string) error {
	command := discord.Command{
		Type:        discord.ChatInputCommand,
		Name:        name,
		Description: description,
	}

	marshalled_command, err := command.MarshalJSON()
	if err != nil {
		return err
	}

	return da.SendApplicationCommandRequest(http.MethodPost, marshalled_command)
}
