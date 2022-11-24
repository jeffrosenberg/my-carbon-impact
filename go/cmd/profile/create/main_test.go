package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/jeffrosenberg/my-carbon-impact/mock/mock_aws"
	"github.com/jeffrosenberg/my-carbon-impact/mock/mock_uuid"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockClient *mock_aws.MockClient
	mockUuid   *mock_uuid.MockUuidGenerator
	ctrl       gomock.Controller
)

const uuidString = "3f198987-b4fd-4594-b2d6-8ac9ae0a44e4"

func initTests(ctrl *gomock.Controller) {
	logger = logging.GetTestLogger()
	mockClient = mock_aws.NewMockClient(ctrl)
	client = mockClient
	mockUuid = mock_uuid.NewMockUuidGenerator(ctrl)
	generator = mockUuid
}

func TestLambdaHandler(t *testing.T) {
	tests := []struct {
		name          string
		skip          bool
		mock          func()
		input         events.APIGatewayProxyRequest
		expected      events.APIGatewayProxyResponse
		expectedError string
	}{
		{
			name: "happy path",
			mock: func() {
				generateId := mockUuid.
					EXPECT().
					NewV7(uuid.MillisecondPrecision).
					Return(uuid.FromStringOrNil(uuidString), nil)
				mockClient.
					EXPECT().
					PutItem(gomock.Any(), gomock.Any()).
					After(generateId).
					Return(&dynamodb.PutItemOutput{}, nil)
			},
			input: events.APIGatewayProxyRequest{
				Body: `{
					"name": "my test",
					"email": "mytest@mailinator.com"
				}`,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body: fmt.Sprintf(`{"id":"%s","name":"%s","email":"%s"}`,
					uuidString, "my test", "mytest@mailinator.com"),
				IsBase64Encoded: false,
			},
		},
		{
			name: "extra JSON fields dropped",
			mock: func() {
				generateId := mockUuid.
					EXPECT().
					NewV7(uuid.MillisecondPrecision).
					Return(uuid.FromStringOrNil(uuidString), nil)
				mockClient.
					EXPECT().
					PutItem(gomock.Any(), gomock.Any()).
					After(generateId).
					Return(&dynamodb.PutItemOutput{}, nil)
			},
			input: events.APIGatewayProxyRequest{
				Body: `{
					"name": "my test",
					"email": "mytest@mailinator.com",
					"my_first_bad_field": "apples",
					"my_first_bad_field": "bananas"
				}`,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body: fmt.Sprintf(`{"id":"%s","name":"%s","email":"%s"}`,
					uuidString, "my test", "mytest@mailinator.com"),
				IsBase64Encoded: false,
			},
		},
		{
			name: "invalid input throws 400 error",
			mock: func() {},
			input: events.APIGatewayProxyRequest{
				Body: `[
					"name": "my test",
					"email": "mytest@mailinator.com"
				]`,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				IsBase64Encoded: false,
			},
			expectedError: "bad request",
		},
		{
			name: "database error",
			mock: func() {
				generateId := mockUuid.
					EXPECT().
					NewV7(uuid.MillisecondPrecision).
					Return(uuid.FromStringOrNil(uuidString), nil)
				mockClient.
					EXPECT().
					PutItem(gomock.Any(), gomock.Any()).
					After(generateId).
					Return(nil, errors.New("database error"))
			},
			input: events.APIGatewayProxyRequest{
				Body: `{
					"name": "my test",
					"email": "mytest@mailinator.com"
				}`,
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusInternalServerError,
				IsBase64Encoded: false,
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

			test.mock()
			got, err := handler(context.Background(), test.input)

			if test.expectedError != "" {
				assert.ErrorContains(t, err, test.expectedError)
			} else if test.expectedError == "" {
				assert.NoError(t, err, "No error expected but received %v", err)
			}
			if test.expected.StatusCode != 0 { // expected is set
				assert.Equal(t, test.expected, got)
			}
		})
	}
}
