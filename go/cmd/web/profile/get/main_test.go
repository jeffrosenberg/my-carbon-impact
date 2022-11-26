package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jeffrosenberg/my-carbon-impact/mock/mock_aws"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock" // for easier assertions
	"github.com/stretchr/testify/assert"
)

var (
	mockClient *mock_aws.MockClient
	ctrl       gomock.Controller
	testUuid   uuid.UUID
)

func initTests(ctrl *gomock.Controller) {
	logger = logging.GetTestLogger()
	mockClient = mock_aws.NewMockClient(ctrl)
	client = mockClient
}

func getTestUuid() uuid.UUID {
	if testUuid.IsNil() {
		gen := uuid.NewGen()
		testUuid, _ = gen.NewV7(uuid.MillisecondPrecision)
	}
	return testUuid
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
				mockClient.
					EXPECT().
					GetItem(gomock.Any(), gomock.Any()). // TODO: Replace gomock.Any() with real expectations
					Return(&dynamodb.GetItemOutput{
						Item: map[string]types.AttributeValue{
							"id":    &types.AttributeValueMemberB{Value: getTestUuid().Bytes()},
							"name":  &types.AttributeValueMemberS{Value: "Jeff"},
							"email": &types.AttributeValueMemberS{Value: "jeff@mailinator.com"},
						},
					}, nil)
			},
			input: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": getTestUuid().String()},
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body: strings.ReplaceAll(fmt.Sprintf(
					`<!DOCTYPE html>
<html>
	<head>
		<title>Test Title</title>
	</head>
	<body>
		<div class="navbar">...</div>
		<div class="content">
<table>
	<thead>
		Profile page: Test Message
	</thead>
	<tbody>
		<tr>
			<td>ID</td>
			<td>%s</td>
		</tr>
		<tr>
			<td>Name</td>
			<td>Jeff</td>
		</tr>
		<tr>
			<td>Email</td>
			<td>jeff@mailinator.com</td>
		</tr>
	</tbody>
</table>
		</div>
	</body>
</html>`, getTestUuid().String()), "\t", "  "),
				IsBase64Encoded: false,
			},
		},
		{
			name:  "missing path parameter throws 400 error",
			mock:  func() {},
			input: events.APIGatewayProxyRequest{},
			expected: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				IsBase64Encoded: false,
			},
			expectedError: "bad request",
		},
		{
			name: "incorrect path parameter throws 400 error",
			mock: func() {},
			input: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"fail": getTestUuid().String()},
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				IsBase64Encoded: false,
			},
			expectedError: "bad request",
		},
		{
			name: "not found throws 404 error",
			mock: func() {
				mockClient.
					EXPECT().
					GetItem(gomock.Any(), gomock.Any()). // TODO: Replace gomock.Any() with real expectations
					Return(&dynamodb.GetItemOutput{}, nil)
			},
			input: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": getTestUuid().String()},
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode:      http.StatusNotFound,
				IsBase64Encoded: false,
			},
			expectedError: "not found",
		},
		{
			name: "database error",
			mock: func() {
				mockClient.
					EXPECT().
					GetItem(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))
			},
			input: events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": getTestUuid().String()},
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
