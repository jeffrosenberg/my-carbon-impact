package interfaces

// TODO: This package is meant to be temporary,
// ideally I'd like to move these somewhere more meaningful

type CarbonCalculator interface {
	Calculate() (int32, error)
}
