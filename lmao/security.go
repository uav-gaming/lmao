package lmao

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/sirupsen/logrus"
)

var discord_public_key []byte

func init() {
	public_key := GetenvMust("DISCORD_PUBLIC_KEY")
	decoded_public_key, err := hex.DecodeString(public_key)
	if err != nil {
		logrus.Fatal("Invalid discord public key: ", public_key)
	}
	discord_public_key = decoded_public_key
}

type VerificationError struct {
}

func VerifyRequest(request Request) bool {
	signature := request.Headers["X-Signature-Ed25519"]
	decoded_signature, err := hex.DecodeString(signature)
	if err != nil {
		logrus.Warn("Invliad request signature: ", signature)
		return false
	}

	timestamp := request.Headers["X-Signature-Timestamp"]
	body := request.Body

	return ed25519.Verify(discord_public_key, []byte(timestamp+body), decoded_signature)
}
