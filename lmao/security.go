package lmao

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/sirupsen/logrus"
)

var discord_public_key []byte

func init() {
	discord_public_key = GetenvMustHex("DISCORD_PUBLIC_KEY")
}

type VerificationError struct {
}

func VerifyRequest(request Request) bool {
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

	return ed25519.Verify(discord_public_key, []byte(timestamp+body), decoded_signature)
}
