package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jeffrosenberg/my-carbon-impact/internal/profile"
	"github.com/stretchr/testify/assert"
)

func initTests() {
	logger = log.With().Logger().Level(zerolog.Disabled)
}

func TestLambdaHandler(t *testing.T) {
	initTests()
	tests := []struct {
		name     string
		input    profile.ProfileInput
		expected events.APIGatewayProxyResponse
	}{
		{
			name:  "happy path",
			input: profile.ProfileInput{},
			expected: events.APIGatewayProxyResponse{
				StatusCode: http.StatusCreated,
			},
		},
	}

	for _, test := range tests {
		event := events.APIGatewayProxyRequest{}
		got, err := handler(context.Background(), event)

		assert.NoErrorf(t, err, "No error expected but received %v", err)
		// Test each field individually instead of comparing the expected object,
		// because we can't know (and don't care) what the generated ID will be
		assert.Equal(t, test.expected, got)
	}
}
