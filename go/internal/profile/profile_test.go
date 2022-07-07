package profile

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProfile(t *testing.T) {
	expected := Profile{
		Name:         "New User",
		Vehicles:     map[string]Vehicle{},
		CarbonEvents: []interfaces.CarbonCalculator{},
	}
	got, err := NewProfile()

	require.NoErrorf(t, err, "No error expected but received %v", err)
	// Test each field individually instead of comparing the expected object,
	// because we can't know (and don't care) what the generated ID will be
	assert.Equal(t, expected.Name, got.Name)
	assert.Equal(t, expected.Vehicles, got.Vehicles)
	assert.Equal(t, expected.CarbonEvents, got.CarbonEvents)
	assert.IsType(t, uuid.UUID{}, got.ID)
}

func TestNewProfileWithInputs(t *testing.T) {
	tests := []struct {
		name  string
		input ProfileInput
	}{
		{
			name: "happy path",
			input: ProfileInput{
				Name: "Philip J. Fry",
				Vehicles: map[string]Vehicle{
					"Planet Express Ship": {
						Year: 3001,
						MPG:  20000,
					},
				},
			},
		},
	}

	for _, test := range tests {
		got, err := NewProfileFromInput(test.input)

		require.NoErrorf(t, err, "No error expected but received %v", err)
		// Test each field individually instead of comparing the expected object,
		// because we can't know (and don't care) what the generated ID will be
		assert.Equal(t, test.input.Name, got.Name)
		assert.Equal(t, test.input.Vehicles, got.Vehicles)
		assert.Equal(t, []interfaces.CarbonCalculator{}, got.CarbonEvents)
		assert.IsType(t, uuid.UUID{}, got.ID)
	}
}
