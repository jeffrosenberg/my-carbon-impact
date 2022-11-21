//go:generate mockgen -destination=../../mock/mock_aws/dynamo.go -package mock_aws . Client,PutItemInputGenerator

package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type PutItemInputGenerator interface {
	GeneratePutItemInput() (*dynamodb.PutItemInput, error)
}
