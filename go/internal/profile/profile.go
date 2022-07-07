package profile

import (
	"github.com/gofrs/uuid"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/interfaces"
)

type Profile struct {
	ID           uuid.UUID
	Name         string
	Vehicles     map[string]Vehicle
	CarbonEvents []interfaces.CarbonCalculator
}

type ProfileInput struct {
	Name     string             `json:"name" validate:"required"`
	Vehicles map[string]Vehicle `json:"vehicles"`
}

func NewProfile() (*Profile, error) {
	id, err := uuid.NewV7(uuid.MillisecondPrecision)
	if err != nil {
		return nil, err
	}

	return &Profile{
		ID:           id,
		Name:         "New User",
		Vehicles:     make(map[string]Vehicle),
		CarbonEvents: make([]interfaces.CarbonCalculator, 0),
	}, nil
}

func NewProfileFromInput(input ProfileInput) (*Profile, error) {
	id, err := uuid.NewV7(uuid.MillisecondPrecision)
	if err != nil {
		return nil, err
	}

	return &Profile{
		ID:           id,
		Name:         input.Name,
		Vehicles:     input.Vehicles,
		CarbonEvents: make([]interfaces.CarbonCalculator, 0),
	}, nil
}
