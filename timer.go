package clock

import (
	"time"
)

type Timer interface {
	Ch() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

type mockTimer struct {
	clock       *MockClock
	id          int
	ch          chan time.Time
	triggerTime time.Time
	isActive    bool

	isTicker bool
	duration time.Duration
}

func (t *mockTimer) Ch() <-chan time.Time {
	return t.ch
}

func (t *mockTimer) Stop() bool {
	return t.clock.stopTimer(t)
}

func (t *mockTimer) Reset(d time.Duration) bool {
	return t.clock.resetTimer(t, d)
}
