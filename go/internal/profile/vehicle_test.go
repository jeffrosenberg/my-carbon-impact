package profile

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVehicleSort(t *testing.T) {
	tests := []struct {
		name     string
		vehicles []OdometerReading
		expected []OdometerReading
	}{
		{
			name: "Simple sort",
			vehicles: []OdometerReading{
				{At: time.Date(2022, 2, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 3, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
			},
			expected: []OdometerReading{
				{At: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 2, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 3, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
			},
		},
		{
			name: "Already sorted",
			vehicles: []OdometerReading{
				{At: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 2, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 3, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
			},
			expected: []OdometerReading{
				{At: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 2, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 3, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
			},
		},
		{
			name: "Nanosecond difference",
			vehicles: []OdometerReading{
				{At: time.Date(2022, 1, 1, 1, 1, 1, 3, time.UTC), Miles: 100},
				{At: time.Date(2022, 1, 1, 1, 1, 1, 2, time.UTC), Miles: 100},
				{At: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
			},
			expected: []OdometerReading{
				{At: time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC), Miles: 100},
				{At: time.Date(2022, 1, 1, 1, 1, 1, 2, time.UTC), Miles: 100},
				{At: time.Date(2022, 1, 1, 1, 1, 1, 3, time.UTC), Miles: 100},
			},
		},
		{
			name:     "Empty slice",
			vehicles: []OdometerReading{},
			expected: []OdometerReading{},
		},
	}

	for _, test := range tests {
		t.Log(test.name)
		sort.Sort(ByDate(test.vehicles))
		assert.Equal(t, test.expected, test.vehicles)
	}
}
