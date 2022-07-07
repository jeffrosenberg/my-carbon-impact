package profile

import (
	"time"
)

type OdometerReading struct {
	At    time.Time
	Miles uint16
}

type ByDate []OdometerReading

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].At.Before(a[j].At) }

type Vehicle struct {
	Year     int               `json:"year"`
	Make     string            `json:"make"`
	Model    string            `json:"model"`
	Odometer []OdometerReading `json:"odometer"`
	MPG      float64           `json:"mpg" validate:"required"` // float64 type simplifies use of math.Round()
}
