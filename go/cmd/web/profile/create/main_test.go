package main

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	ctrl gomock.Controller
)

func initTests(ctrl *gomock.Controller) {
	logger = logging.GetTestLogger()
}

func TestLambdaHandler(t *testing.T) {
	tests := []struct {
		name          string
		skip          bool
		input         events.APIGatewayProxyRequest
		expected      events.APIGatewayProxyResponse
		expectedError string
	}{
		{
			name: "happy path",
			input: events.APIGatewayProxyRequest{
				Resource:   "/web/profile",
				Path:       "",
				HTTPMethod: "GET",
			},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body: strings.ReplaceAll(
					`<!DOCTYPE html>
<html>
	<head>
		<title>Test Title</title>
	</head>
	<body>
		<div class="navbar">...</div>
		<div class="content">
<div>
	<p>This is just a placeholder for now</p>
	<p>Test Message</p>
</div>
		</div>
	</body>
</html>`, "\t", "  "),
				IsBase64Encoded: false,
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
