package lmao

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Get the environment variable.
// Crash the program if it's not set.
func GetenvMust(name string) string {
	value := os.Getenv(name)
	if len(value) <= 0 {
		logrus.Fatal("env ", name, " not set")
	}
	return value
}
