package calculators

import (
	"testing"
	"time"

	"github.com/jeffrosenberg/my-carbon-impact/internal/profile"
	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	tests := []struct {
		name           string
		skip           bool
		input          VehicleMileageInput
		expected       VehicleMileageEvent
		expectedErrors []string
	}{
		{
			name: "Happy path",
			input: VehicleMileageInput{
				Vehicle: &profile.Vehicle{
					MPG: 20,
					Odometer: []profile.OdometerReading{
						{
							At:    time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
							Miles: 900,
						},
					},
				},
				At:    "2022-07-01T00:00:00Z",
				Miles: 1000,
			},
			expected: VehicleMileageEvent{
				From:  time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
				Until: time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC),
				Miles: 100,
				MPG:   20,
			},
		},
		{
			name: "Converts inputs to UTC",
			input: VehicleMileageInput{
				Vehicle: &profile.Vehicle{
					MPG: 20,
					Odometer: []profile.OdometerReading{
						{
							At:    time.Date(2022, 7, 1, 5, 0, 0, 0, time.UTC),
							Miles: 900,
						},
					},
				},
				At:    "2022-08-01T00:00:00-05:00",
				Miles: 1000,
			},
			expected: VehicleMileageEvent{
				From:  time.Date(2022, 7, 1, 5, 0, 0, 0, time.UTC),
				Until: time.Date(2022, 8, 1, 5, 0, 0, 0, time.UTC),
				Miles: 100,
				MPG:   20,
			},
		},
		{
			name: "Validation error: Vehicle",
			input: VehicleMileageInput{
				At:    "2022-08-01T00:00:00-05:00",
				Miles: 1000,
			},
			expectedErrors: []string{"'Vehicle'"},
		},
		{
			name: "Validation error: At",
			input: VehicleMileageInput{
				Vehicle: &profile.Vehicle{},
				Miles:   1000,
			},
			expectedErrors: []string{"'At'"},
		},
		{
			name: "Validation error: Miles",
			input: VehicleMileageInput{
				Vehicle: &profile.Vehicle{},
				Miles:   1000,
			},
			expectedErrors: []string{"'At'"},
		},
		{
			name:           "Empty input",
			input:          VehicleMileageInput{},
			expectedErrors: []string{"'Vehicle'", "'At'", "'Miles'"},
		},
	}

	for _, test := range tests {
		if test.skip {
			t.Skipf("Skipping %s", test.name)
		}
		t.Log(test.name)

		got, err := New(test.input)
		if len(test.expectedErrors) > 0 {
			assert.Error(t, err, "Expected an error but none occurred")
			for _, expected := range test.expectedErrors {
				assert.ErrorContains(t, err, expected)
			}
		} else {
			assert.NoErrorf(t, err, "No error expected but got %w", err)
			assert.Equal(t, test.expected, got)
		}
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		testEvent VehicleMileageEvent
		expected  int32
	}{
		{
			name: "simple calculation",
			testEvent: VehicleMileageEvent{
				MPG:   20,
				Miles: 100,
			},
			expected: 99, // 5 * 19.6 * 1.01 = 98.98,
		},
		{
			name: "round up",
			testEvent: VehicleMileageEvent{
				MPG:   20,
				Miles: 99,
			},
			expected: 98, // (99/20) * 19.6 * 1.01 = 97.9902,
		},
		{
			name: "round down",
			testEvent: VehicleMileageEvent{
				MPG:   20,
				Miles: 98,
			},
			expected: 97, // 5 * 19.6 * 1.01 = 97.0004,
		},
		{
			name: "small event",
			testEvent: VehicleMileageEvent{
				MPG:   30,
				Miles: 2,
			},
			expected: 1,
		},
		{
			name: "zero miles", // Just returning zero for now. TODO: would we prefer an error?
			testEvent: VehicleMileageEvent{
				MPG:   20,
				Miles: 0,
			},
			expected: 0,
		},
		{
			name: "large event",
			testEvent: VehicleMileageEvent{
				MPG:   18,
				Miles: 990,
			},
			expected: 1089,
		},
	}

	for _, test := range tests {
		t.Log(test.name)
		got := test.testEvent.Calculate()
		// assert.NoErrorf(t, err, "No error expected but got %w", err)
		assert.Equal(t, test.expected, got)
	}
}

func TestTimeZoneLocalization(t *testing.T) {
	t.Skip("TODO: build and test time zone localization")
}
