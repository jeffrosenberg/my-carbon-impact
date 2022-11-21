//go:generate mockgen -destination=../../mock/mock_uuid/uuid.go -package mock_uuid . UuidGenerator

package interfaces

// TODO: This package is meant to be temporary,
// ideally I'd like to move these somewhere more meaningful

import (
	"github.com/gofrs/uuid"
)

type CarbonCalculator interface {
	Calculate() (int32, error)
}

type UuidGenerator interface {
	uuid.Generator
}
