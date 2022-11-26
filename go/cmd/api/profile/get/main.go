package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofrs/uuid"

	"github.com/jeffrosenberg/my-carbon-impact/internal/db"
	"github.com/jeffrosenberg/my-carbon-impact/internal/db/get"
	"github.com/jeffrosenberg/my-carbon-impact/internal/profile"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"
	"github.com/rs/zerolog"
)

var (
	client db.Client
	logger zerolog.Logger
)

func initFunc() {
	logger = logging.GetLogger().
		With().
		Str("entity", "profile").
		Str("operation", "get").
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
}

func main() {
	initFunc()
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (event events.APIGatewayProxyResponse, err error) {
	logger = logging.AppendContext(&logger, ctx).With().Logger()
	logger.Trace().Msg("Get profile lambda beginning")
	logger = logger.With().Interface("e", e).Logger()

	id := e.PathParameters["id"]
	if id == "" {
		logger.Warn().Msg("No id parameter sent")
		event.StatusCode = http.StatusBadRequest
		return event, errors.New("bad request: missing profile id")
	}

	input := &get.DynamoDbGetInput{
		Id:     uuid.FromStringOrNil(id),
		Ctx:    ctx,
		Client: client,
		Logger: logger,
	}
	res, err := get.GetProfile(input)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read profile from DynamoDb")
		event.StatusCode = http.StatusInternalServerError
		return
	}
	logger.Debug().Interface("res", res).Msg("Got profile from DynamoDb")

	if res.Item == nil { // no records found
		event.StatusCode = http.StatusNotFound
		return event, fmt.Errorf("not found: no profile found with id %s", id)
	}

	var profile profile.Profile
	err = attributevalue.UnmarshalMap(res.Item, &profile)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to Unmarshall profile")
		event.StatusCode = http.StatusInternalServerError
		return
	}

	jsonBytes, err := json.Marshal(profile)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshall profile to JSON")
		event.StatusCode = http.StatusInternalServerError
		err = nil
		return
	}

	event.StatusCode = http.StatusOK
	event.Body = string(jsonBytes)
	logger.Info().Interface("profile", profile).Msg("Get profile lambda complete")
	return
}
