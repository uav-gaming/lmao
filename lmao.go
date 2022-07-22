package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/sirupsen/logrus"
	"github.com/tjhu/lmao/lmao"
)

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

	response, err := lmao.HandleEvent(event)
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

	lambda.Start(HandleRequest)
}
