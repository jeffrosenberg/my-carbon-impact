package get

import (
	"context"

	"github.com/jeffrosenberg/my-carbon-impact/internal/db"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/constants"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
)

type DynamoDbGetInput struct {
	Id     uuid.UUID       `json:"id"`
	Ctx    context.Context `json:"-"`
	Client db.Client       `json:"-"`
	Logger zerolog.Logger  `json:"-"`
}

func GetProfile(req *DynamoDbGetInput) (*dynamodb.GetItemOutput, error) {
	req.Logger.Trace().Msg("Get profile")

	input := &dynamodb.GetItemInput{
		TableName: &constants.DYNAMO_TABLE_NAME,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberB{Value: req.Id.Bytes()},
		},
	}

	output, err := req.Client.GetItem(req.Ctx, input)
	if err != nil {
		req.Logger.Error().Err(err).Msg("Failed to get DynamoDb item")
		return nil, err
	}

	return output, nil
}
