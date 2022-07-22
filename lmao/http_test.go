package lmao

import (
	"encoding/json"
	"testing"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/stretchr/testify/assert"
)

const REQUEST string = `{
	"version": "2.0",
	"routeKey": "$default",
	"rawPath": "/",
	"rawQueryString": "",
	"headers": {
	  "content-length": "123",
	  "x-amzn-tls-cipher-suite": "ECDHE-RSA-AES128-GCM-SHA256",
	  "x-signature-ed25519": "abc123",
	  "x-amzn-tls-version": "TLSv1.2",
	  "x-amzn-trace-id": "Root=1-abc123-abc123",
	  "x-forwarded-proto": "https",
	  "host": "xyz.lambda-url.us-east-1.on.aws",
	  "x-forwarded-port": "443",
	  "content-type": "application/json",
	  "x-forwarded-for": "1.1.1.1",
	  "x-signature-timestamp": "111",
	  "user-agent": "Discord-Interactions/1.0 (+https://discord.com)"
	},
	"requestContext": {
	  "accountId": "anonymous",
	  "apiId": "xyz",
	  "domainName": "xyz.lambda-url.us-east-1.on.aws",
	  "domainPrefix": "xyz",
	  "http": {
		"method": "POST",
		"path": "/",
		"protocol": "HTTP/1.1",
		"sourceIp": "1.1.1.1",
		"userAgent": "Discord-Interactions/1.0 (+https://discord.com)"
	  },
	  "requestId": "1a-2b-3c-4d-5e",
	  "routeKey": "$default",
	  "stage": "$default",
	  "time": "09/Jun/2022:11:22:33 +0000",
	  "timeEpoch": 123
	},
	"body": "{\"application_id\":\"123\",\"id\":\"123\",\"token\":\"xyz\",\"type\":1,\"user\":{\"avatar\":\"xyz\",\"avatar_decoration\":null,\"discriminator\":\"1234\",\"id\":\"123\",\"public_flags\":0,\"username\":\"dude\"},\"version\":1}",
	"isBase64Encoded": false
  }`

func TestRequestUnmarshal(t *testing.T) {
	var request Request
	assert.NoError(t, json.Unmarshal([]byte(REQUEST), &request))
	assert.Equal(t, "POST", request.RequestContext.HTTP.Method)

	var event discord.InteractionEvent
	assert.NoError(t, event.UnmarshalJSON([]byte(request.Body)))
	assert.Equal(t, "dude", event.User.Username)
}
