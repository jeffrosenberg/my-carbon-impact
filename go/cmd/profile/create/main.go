package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofrs/uuid"

	"github.com/jeffrosenberg/my-carbon-impact/internal/db"
	"github.com/jeffrosenberg/my-carbon-impact/internal/db/create"
	"github.com/jeffrosenberg/my-carbon-impact/internal/profile"
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

	var profileInput profile.ProfileInput
	err = json.Unmarshal([]byte(e.Body), &profileInput)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshall request body")
		event.StatusCode = http.StatusBadRequest
		return event, fmt.Errorf("bad request: %w", err)
	}

	profile, err := profile.NewProfileFromInput(profileInput, generator)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create profile from input")
		event.StatusCode = http.StatusInternalServerError
		return
	}

	createProfileRequest := create.DynamoDbCreateInput{
		Input:  profile,
		Ctx:    ctx,
		Client: client,
		Logger: logger,
	}
	logger.Debug().Interface("create_profile_request", createProfileRequest).Msg("Creating profile in DynamoDb")
	_, err = create.CreateProfile(createProfileRequest)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write profile to DynamoDb")
		event.StatusCode = http.StatusInternalServerError
		return
	}

	jsonBytes, err := json.Marshal(profile)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshall profile to JSON")
		event.StatusCode = http.StatusCreated
		err = nil
		return
	}

	event.StatusCode = http.StatusOK
	event.Body = string(jsonBytes)
	logger.Info().Interface("profile", profile).Msg("Create profile lambda complete")
	return
}
