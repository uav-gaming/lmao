package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/tjhu/lmao/lmao"
)

func HandleRequest(ctx context.Context, request lmao.Request) (string, error) {
	log.Info(request)
	return "Hello world!", nil
}

func main() {
	// Setup logger.
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true, // Logged by AWS.
	})
	log.SetReportCaller(true) // Log caller.

	lambda.Start(HandleRequest)
}
