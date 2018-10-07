package clock

import (
	"time"
)

// Timer is an interface that represents both real and mock timers.
type Timer interface {
	Chan() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

var _ Timer = realTimer{}
var _ Timer = mockTimer{}

// realTimer is just a time.Timer's shallow wrapper.
type realTimer struct {
	*time.Timer
}

// Chan implementats Timer.
func (t realTimer) Chan() <-chan time.Time {
	return t.C
}

// mockTimer is just an internalTimer's shallow wrapper.
type mockTimer struct {
	*internalTimer
}

// Chan implementats Timer.
func (t mockTimer) Chan() <-chan time.Time {
	return t.ch
}

// Stop implementats Timer.
func (t mockTimer) Stop() bool {
	return t.clock.stopTimer(t.internalTimer)
}

// Reset implementats Timer.
func (t mockTimer) Reset(d time.Duration) bool {
	return t.clock.resetTimer(t.internalTimer, d)
}
