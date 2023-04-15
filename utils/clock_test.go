package utils

import "time"

var _ Clock = FixedClock{}

type FixedClock struct {
	now time.Time
}

func (f FixedClock) Now() time.Time {
	return f.now
}
