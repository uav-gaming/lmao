package lmao

import (
	"encoding/hex"
	"os"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/sirupsen/logrus"
)

// Get the environment variable.
// Crash the program if it's not set.
func GetenvMustStr(name string) string {
	value := os.Getenv(name)
	if len(value) <= 0 {
		logrus.Fatal("env ", name, " not set")
	}
	return value
}

func GetenvMustHex(name string) []byte {
	value := GetenvMustStr(name)
	decoded_value, err := hex.DecodeString(value)
	if err != nil {
		logrus.Fatal("Invalid discord public key: ", value)
	}
	return decoded_value
}

func GetenvMustUint64(name string) uint64 {
	value, err := strconv.ParseUint(GetenvMustStr(name), 10, 64)
	if err != nil {
		logrus.Fatal("Failed to convert ", value, " to uint64")
	}
	return value
}

func GetenvMustSnowflake(name string) discord.Snowflake {
	return discord.Snowflake(GetenvMustUint64(name))
}
