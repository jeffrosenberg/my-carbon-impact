package calculators

import (
	"math"
	"sort"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jeffrosenberg/my-carbon-impact/internal/profile"
	epa "github.com/jeffrosenberg/my-carbon-impact/pkg/constants"
)

type VehicleMileageInput struct {
	Vehicle *profile.Vehicle `validate:"required"`
	At      string           `validate:"required"`
	Miles   uint16           `validate:"required"`
}

type VehicleMileageEvent struct {
	From       time.Time
	Until      time.Time
	Miles      uint16
	MPG        float64
	CarbonCost int32
}

func (event VehicleMileageEvent) Calculate() int32 {
	return int32(math.Round(
		(float64(event.Miles) / event.MPG) *
			epa.VEHICLE_LBS_PER_GAL *
			epa.VEHICLE_OTHER_EMISSIONS_RATIO))
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// TODO: Also add to the profile's Odometer readings?
func New(input VehicleMileageInput) (event VehicleMileageEvent, err error) {
	err = validate.Struct(input)
	if err != nil {
		return
	}

	at, err := time.Parse(time.RFC3339, input.At)
	if err != nil {
		return
	}

	sort.Sort(profile.ByDate(input.Vehicle.Odometer))

	event = VehicleMileageEvent{
		From:  input.Vehicle.Odometer[len(input.Vehicle.Odometer)-1].At,
		Until: at.UTC(),
		Miles: input.Miles - input.Vehicle.Odometer[len(input.Vehicle.Odometer)-1].Miles,
		MPG:   input.Vehicle.MPG,
	}

	return event, nil
}
