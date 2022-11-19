package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"
	"github.com/rs/zerolog"
)

var (
	logger zerolog.Logger
)

func initLogger(ctx context.Context) {
	logger =
		logging.GetLoggerWithContext(ctx).
			With().
			Str("entity", "profile").
			Str("operation", "list").
			Logger()
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (event events.APIGatewayProxyResponse, err error) {
	initLogger(ctx)

	event = events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            "",
		IsBase64Encoded: false,
	}
	return
}
