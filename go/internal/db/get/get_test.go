package get

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"

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

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name          string
		skip          bool
		mock          func()
		expected      *dynamodb.GetItemOutput
		expectedError string
	}{
		{
			name: "happy path",
			mock: func() {
				mockClient.
					EXPECT().
					GetItem(gomock.Any(), gomock.Any()). // TODO: Replace gomock.Any() with real expectations
					Return(&dynamodb.GetItemOutput{
						Item: map[string]types.AttributeValue{
							"id":    &types.AttributeValueMemberS{Value: getTestUuid().String()},
							"name":  &types.AttributeValueMemberS{Value: "Jeff"},
							"email": &types.AttributeValueMemberS{Value: "jeff@mailinator.com"},
						},
					}, nil)
			},
			expected: &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"id":    &types.AttributeValueMemberS{Value: getTestUuid().String()},
					"name":  &types.AttributeValueMemberS{Value: "Jeff"},
					"email": &types.AttributeValueMemberS{Value: "jeff@mailinator.com"},
				},
			},
		},
		{
			name: "no record found",
			mock: func() {
				mockClient.
					EXPECT().
					GetItem(gomock.Any(), gomock.Any()). // TODO: Replace gomock.Any() with real expectations
					Return(&dynamodb.GetItemOutput{}, nil)
			},
			expected: &dynamodb.GetItemOutput{},
		},
		{
			name: "database error",
			mock: func() {
				mockClient.
					EXPECT().
					GetItem(gomock.Any(), gomock.Any()). // TODO: Replace gomock.Any() with real expectations
					Return(&dynamodb.GetItemOutput{}, errors.New("database error"))
			},
			expectedError: "database error",
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

			req := &DynamoDbGetInput{
				Id:     getTestUuid(),
				Ctx:    context.Background(),
				Client: mockClient,
				Logger: logger,
			}
			test.mock()
			got, err := GetProfile(req)
			if test.expectedError != "" {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err, "No error expected but received %v", err)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}
