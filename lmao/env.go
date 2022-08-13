package lmao

import (
	"encoding/hex"
	"os"
	"reflect"
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

type ISnowflake interface {
	discord.Snowflake | discord.AppID | discord.UserID
	IsValid() bool
}

func GetenvMustValidSnowflake[T ISnowflake](name string) T {
	val := T(GetenvMustUint64(name))
	if !val.IsValid() {
		logrus.Fatal("Invalid ", reflect.TypeOf(val), " value of ", val, " from ", name)
	}
	return val
}
