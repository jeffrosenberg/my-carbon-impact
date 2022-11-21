package profile

import (
	"testing"

	"github.com/jeffrosenberg/my-carbon-impact/mock/mock_uuid"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/constants"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/interfaces"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mockUuid *mock_uuid.MockUuidGenerator
	ctrl     gomock.Controller
	logger   zerolog.Logger
)

const uuidString string = "c03e7835-17ed-421b-8c59-50f3266f71e8"

func initTests(ctrl *gomock.Controller) {
	logger = logging.GetTestLogger()
	mockUuid = mock_uuid.NewMockUuidGenerator(ctrl)
}

func TestNewProfile(t *testing.T) {
	expected := &Profile{
		ID:           uuid.FromStringOrNil(uuidString),
		Name:         "New User",
		Vehicles:     map[string]Vehicle{},
		CarbonEvents: []interfaces.CarbonCalculator{},
	}
	ctrl = *gomock.NewController(t)
	initTests(&ctrl)

	mockUuid.
		EXPECT().
		NewV7(gomock.Any()).
		Return(uuid.FromStringOrNil(uuidString), nil)
	got, err := NewProfile(mockUuid)

	require.NoErrorf(t, err, "No error expected but received %v", err)
	assert.Equal(t, expected, got)
}

func TestNewProfileWithInputs(t *testing.T) {
	input := ProfileInput{
		Name:  "Philip J. Fry",
		Email: "philipjfry@mailinator.com",
		Vehicles: map[string]Vehicle{
			"Planet Express Ship": {
				Year: 3001,
				MPG:  20000,
			},
		},
	}
	expected := &Profile{
		ID:    uuid.FromStringOrNil(uuidString),
		Name:  "Philip J. Fry",
		Email: "philipjfry@mailinator.com",
		Vehicles: map[string]Vehicle{
			"Planet Express Ship": {
				Year: 3001,
				MPG:  20000,
			},
		},
		CarbonEvents: []interfaces.CarbonCalculator{},
	}

	ctrl = *gomock.NewController(t)
	initTests(&ctrl)

	mockUuid.
		EXPECT().
		NewV7(gomock.Any()).
		Return(uuid.FromStringOrNil(uuidString), nil)
	got, err := NewProfileFromInput(input, mockUuid)

	assert.NoErrorf(t, err, "No error expected but received %v", err)
	assert.Equal(t, expected, got)
}

func TestGeneratePutItemInput(t *testing.T) {
	tests := []struct {
		name          string
		skip          bool
		profile       *Profile
		expected      *dynamodb.PutItemInput
		expectedError string
	}{
		{
			name: "happy path",
			profile: &Profile{
				ID:    uuid.FromStringOrNil(uuidString),
				Name:  "Philip J. Fry",
				Email: "philipjfry@mailinator.com",
				Vehicles: map[string]Vehicle{
					"Planet Express Ship": {
						Year: 3001,
						MPG:  20000,
					},
				},
				CarbonEvents: []interfaces.CarbonCalculator{},
			},
			expected: &dynamodb.PutItemInput{
				TableName: &constants.DYNAMO_TABLE_NAME,
				Item: map[string]types.AttributeValue{
					"id":    &types.AttributeValueMemberB{Value: uuid.FromStringOrNil(uuidString).Bytes()},
					"name":  &types.AttributeValueMemberS{Value: "Philip J. Fry"},
					"email": &types.AttributeValueMemberS{Value: "philipjfry@mailinator.com"},
				},
			},
		},
		{
			name: "marshall error",
			skip: true, // TODO: Find a way to actually trigger this error!
			profile: &Profile{
				Name: "Marshall error",
			},
			expectedError: "unable to marshall Profile to PutItemInput",
		},
	}

	for _, test := range tests {
		if test.skip {
			t.Skipf("Skipping %s", test.name)
		}

		t.Run(test.name, func(t *testing.T) {
			t.Log(test.name)

			got, err := test.profile.GeneratePutItemInput()
			if test.expectedError != "" {
				assert.ErrorContains(t, err, test.expectedError)
				assert.ErrorAs(t, err, attributevalue.InvalidMarshalError{})
			} else {
				assert.NoError(t, err, "No error expected but received %v", err)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}
