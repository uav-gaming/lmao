package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/sirupsen/logrus"
	"github.com/tjhu/lmao/lmao"
)

// Constants
const BUILD_TIME_KEY string = "build_time"

// Globals
var (
	// Variables populated by build script.
	BuildTime string

	// Variables populated by `init`s.
	BuildInfo map[string]interface{}
)

// Initialize build info
func init() {
	BuildInfo = make(map[string]interface{})

	// Get build info from binary.
	info, ok := debug.ReadBuildInfo()
	if !ok {
		logrus.Fatal("Failed to read build info.")
	}
	for _, setting := range info.Settings {
		BuildInfo[setting.Key] = setting.Value
	}

	// Parse and inject build time
	build_epoch, err := strconv.ParseInt(BuildTime, 10, 64)
	if err != nil {
		log.Fatal("Failed to parse timestamp to int64: ", err)
	}
	BuildInfo[BUILD_TIME_KEY] = time.Unix(build_epoch, 0)
}

func HandleRequest(ctx context.Context, request lmao.Request) (lmao.Response, error) {
	logrus.Info(request)

	if !lmao.VerifyRequest(request) {
		logrus.Warn("Request verification failed.")
		return lmao.Response{
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	var event discord.InteractionEvent
	if err := event.UnmarshalJSON([]byte(request.Body)); err != nil {
		logrus.Error("Failed to unmarshal request body: ", request.Body)
		return lmao.Response{}, errors.New("invalid request format")
	}

	response, err := lmao.HandleInteraction(event)
	if err != nil {
		return err.ToResponse(), nil
	}
	return lmao.ToHttpResponse(response)
}

func main() {
	// Setup logger.
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true, // Logged by AWS.
	})
	logrus.SetReportCaller(true) // Log caller.

	logrus.Info("Startup with build info: ", BuildInfo)
	lambda.Start(HandleRequest)
}
