package discordapi

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// A threadsafe abstraction for making requests to discord REST APIs.
type DiscordApi struct {
	// Endpoint for sending global requests.
	// Format: "`DISCORD_APPLICATION_URL`/<application_id>"
	application_url string
	// For the "Authorization" in the request header.
	// Format: "Bot <application_token>"
	authorization string
}

func NewDiscordApi(application_id string, token string) *DiscordApi {
	return &DiscordApi{
		application_url: fmt.Sprintf("%s/%s", DISCORD_APPLICATION_URL, application_id),
		authorization:   fmt.Sprintf("Bot %s", token),
	}
}

// Send request to "`endpoint`" with `body` as the json body.
func (da *DiscordApi) SendRequest(method string, endpoint string, body []byte) error {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", da.authorization)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	logrus.Info(res)
	return nil
}
