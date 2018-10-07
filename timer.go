package clock

import (
	"time"
)

type Timer interface {
	Chan() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

var _ Timer = realTimer{}
var _ Timer = mockTimer{}

type realTimer struct {
	*time.Timer
}

func (t realTimer) Chan() <-chan time.Time {
	return t.C
}

type mockTimer struct {
	*internalTimer
}

func (t mockTimer) Chan() <-chan time.Time {
	return t.ch
}

func (t mockTimer) Stop() bool {
	return t.clock.stopTimer(t.internalTimer)
}

func (t mockTimer) Reset(d time.Duration) bool {
	return t.clock.resetTimer(t.internalTimer, d)
}
