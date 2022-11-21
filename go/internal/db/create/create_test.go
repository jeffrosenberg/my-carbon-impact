package create

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jeffrosenberg/my-carbon-impact/internal/profile"

	"github.com/jeffrosenberg/my-carbon-impact/mock/mock_aws"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	mockClient *mock_aws.MockClient
	ctrl       gomock.Controller
	logger     zerolog.Logger
	testUuid   uuid.UUID
)

func initTests(ctrl *gomock.Controller) {
	logger = logging.GetTestLogger()
	mockClient = mock_aws.NewMockClient(ctrl)
}

func getTestUuid() uuid.UUID {
	if testUuid.IsNil() {
		gen := uuid.NewGen()
		testUuid, _ = gen.NewV7(uuid.MillisecondPrecision)
	}
	return testUuid
}

func TestCreateProfile(t *testing.T) {
	tests := []struct {
		name          string
		skip          bool
		input         *profile.Profile
		mock          func()
		expected      *dynamodb.PutItemOutput
		expectedError string
	}{
		{
			name: "happy path",
			input: &profile.Profile{
				ID:    getTestUuid(),
				Name:  "Jeff",
				Email: "jeff@mailinator.com",
			},
			mock: func() {
				mockClient.
					EXPECT().
					PutItem(gomock.Any(), gomock.Any()).
					Return(&dynamodb.PutItemOutput{
						Attributes: map[string]types.AttributeValue{
							"id":    &types.AttributeValueMemberS{Value: getTestUuid().String()},
							"name":  &types.AttributeValueMemberS{Value: "Jeff"},
							"email": &types.AttributeValueMemberS{Value: "jeff@mailinator.com"},
						},
					}, nil)
			},
			expected: &dynamodb.PutItemOutput{
				Attributes: map[string]types.AttributeValue{
					"id":    &types.AttributeValueMemberS{Value: getTestUuid().String()},
					"name":  &types.AttributeValueMemberS{Value: "Jeff"},
					"email": &types.AttributeValueMemberS{Value: "jeff@mailinator.com"},
				},
			},
		},
	}

	for _, test := range tests {
		if test.skip {
			t.Skipf("Skipping %s", test.name)
		}

		t.Run(test.name, func(t *testing.T) {
			t.Log(test.name)
			ctrl = *gomock.NewController(t)
			initTests(&ctrl)

			req := DynamoDbCreateInput{
				Input:  test.input,
				Ctx:    context.Background(),
				Client: mockClient,
				Logger: logger,
			}
			test.mock()
			got, err := CreateProfile(req)
			if test.expectedError != "" {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err, "No error expected but received %v", err)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}
