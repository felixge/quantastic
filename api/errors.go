package api

import (
	"fmt"
	"time"
)

func NewTimeoutError(d time.Duration) *TimeoutError {
	return &TimeoutError{d}
}

type TimeoutError struct {
	Duration time.Duration
}

func (t *TimeoutError) Error() string {
	return fmt.Sprintf("Timeout exceeded: %s", t.Duration)
}
