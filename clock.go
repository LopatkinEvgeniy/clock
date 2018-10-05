package clock

import (
	"time"
)

// TODO: After
// TODO: AfterFunc (specify callback to timers)
// TODO: Sleep
// TODO: Since
// TODO: Until
// TODO: Tick

type MockClock struct {
	*internalClock
}

func NewMock() *MockClock {
	return &MockClock{
		internalClock: newInternalClock(),
	}
}

func (c *MockClock) Now() time.Time {
	return c.getCurrentTime()
}

func (c *MockClock) Add(d time.Duration) {
	c.moveTimeForward(d)
}

func (c *MockClock) NewTicker(d time.Duration) Ticker {
	return c.makeMockTimer(d, true)
}

func (c *MockClock) NewTimer(d time.Duration) Timer {
	return c.makeMockTimer(d, false)
}
