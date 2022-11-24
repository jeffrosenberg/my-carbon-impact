package create

import (
	"context"

	"github.com/jeffrosenberg/my-carbon-impact/internal/db"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog"
)

type DynamoDbCreateInput struct {
	Input  db.PutItemInputGenerator `json:"input,omitempty"`
	Ctx    context.Context          `json:"-"`
	Client db.Client                `json:"-"`
	Logger zerolog.Logger           `json:"-"`
}

func CreateProfile(req *DynamoDbCreateInput) (*dynamodb.PutItemOutput, error) {
	req.Logger.Trace().Msg("Create profile")

	dbInput, err := req.Input.GeneratePutItemInput()
	if err != nil {
		req.Logger.Error().Err(err).Msg("Failed to generate PutItemInput")
		return nil, err
	}

	output, err := req.Client.PutItem(req.Ctx, dbInput)
	if err != nil {
		req.Logger.Error().Err(err).Msg("Failed to create DynamoDb item")
		return nil, err
	}

	return output, nil
}
