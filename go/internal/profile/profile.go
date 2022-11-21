package profile

import (
	"fmt"

	"github.com/jeffrosenberg/my-carbon-impact/pkg/constants"
	itf "github.com/jeffrosenberg/my-carbon-impact/pkg/interfaces"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofrs/uuid"
)

type Profile struct {
	ID           uuid.UUID              `json:"id" dynamodbav:"id"`
	Name         string                 `json:"name,omitempty" dynamodbav:"name"`
	Email        string                 `json:"email" dynamodbav:"email"`
	Vehicles     map[string]Vehicle     `json:"vehicles,omitempty" dynamodbav:"-"`
	CarbonEvents []itf.CarbonCalculator `json:"-" dynamodbav:"-"`
}

type ProfileInput struct {
	Name     string             `json:"name" validate:"required"`
	Email    string             `json:"email" validate:"required"`
	Vehicles map[string]Vehicle `json:"vehicles"`
}

func NewProfile(idgen itf.UuidGenerator) (*Profile, error) {
	id, err := idgen.NewV7(uuid.MillisecondPrecision)
	if err != nil {
		return nil, err
	}

	return &Profile{
		ID:           id,
		Name:         "New User",
		Vehicles:     make(map[string]Vehicle),
		CarbonEvents: make([]itf.CarbonCalculator, 0),
	}, nil
}

func NewProfileFromInput(input ProfileInput, idgen itf.UuidGenerator) (*Profile, error) {
	id, err := idgen.NewV7(uuid.MillisecondPrecision)
	if err != nil {
		return nil, err
	}

	return &Profile{
		ID:           id,
		Name:         input.Name,
		Email:        input.Email,
		Vehicles:     input.Vehicles,
		CarbonEvents: make([]itf.CarbonCalculator, 0),
	}, nil
}

// fulfills db.PutItemInputGenerator interface
func (p *Profile) GeneratePutItemInput() (*dynamodb.PutItemInput, error) {
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall Profile to PutItemInput: %w", err)
	}
	return &dynamodb.PutItemInput{
		TableName: &constants.DYNAMO_TABLE_NAME,
		Item:      item,
	}, nil
}
