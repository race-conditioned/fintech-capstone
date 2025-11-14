package limiter

import "time"

// Clock allows tests to control time deterministically.
type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() } // includes monotonic component
