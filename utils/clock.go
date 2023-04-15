package utils

import "time"

var _ Clock = UTCClock{}

type Clock interface {
	Now() time.Time
}

type UTCClock struct{}

func (c UTCClock) Now() time.Time {
	return time.Now().UTC()
}
