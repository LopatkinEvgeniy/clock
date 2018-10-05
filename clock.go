package clock

import (
	"time"
)

type MockClock struct {
	*internalClock
}

func NewMock() *MockClock {
	return &MockClock{
		internalClock: newInternalClock(),
	}
}

func (c *MockClock) Add(d time.Duration) {
	c.moveTimeForward(d)
}

func (c *MockClock) Now() time.Time {
	return c.getCurrentTime()
}

func (c *MockClock) After(d time.Duration) <-chan time.Time {
	return c.NewTimer(d).Ch()
}

func (c *MockClock) AfterFunc(d time.Duration, f func()) Timer {
	return c.makeMockTimer(d, false, f)
}

func (c *MockClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

func (c *MockClock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}

func (c *MockClock) Sleep(d time.Duration) {
	<-c.NewTimer(d).Ch()
}

func (c *MockClock) Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return c.NewTicker(d).Ch()
}

func (c *MockClock) NewTicker(d time.Duration) Ticker {
	return c.makeMockTimer(d, true, nil)
}

func (c *MockClock) NewTimer(d time.Duration) Timer {
	return c.makeMockTimer(d, false, nil)
}
