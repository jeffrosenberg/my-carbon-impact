package main

import (
	"bytes"
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofrs/uuid"

	"github.com/jeffrosenberg/my-carbon-impact/internal/db"
	"github.com/jeffrosenberg/my-carbon-impact/internal/html"
	itf "github.com/jeffrosenberg/my-carbon-impact/pkg/interfaces"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"
	"github.com/rs/zerolog"
)

var (
	client    db.Client
	generator itf.UuidGenerator
	logger    zerolog.Logger
)

func initFunc() {
	logger = logging.GetLogger().
		With().
		Str("entity", "profile").
		Str("operation", "create").
		Str("tier", "web").
		Logger()

	region := os.Getenv("region")
	if region == "" {
		region = "us-west-2"
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to load AWS SDK")
	}
	client = dynamodb.NewFromConfig(cfg)

	generator = uuid.NewGen()
}

func main() {
	initFunc()
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (event events.APIGatewayProxyResponse, err error) {
	logger = logging.AppendContext(&logger, ctx).With().Logger()
	logger.Trace().Msg("Create profile lambda beginning")
	logger = logger.With().Interface("e", e).Logger()

	var w bytes.Buffer
	params := html.ProfileParams{
		Title:   "Test Title",
		Message: "Test Message",
	}
	err = html.CreateProfile(&w, params)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to render profile page")
		event.StatusCode = http.StatusInternalServerError
		return
	}

	event.StatusCode = http.StatusOK
	event.Body = w.String()
	logger.Info().Msg("Get profile lambda complete")
	return
}
