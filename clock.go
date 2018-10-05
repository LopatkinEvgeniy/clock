package clock

import (
	"time"
)

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) Timer
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	NewTicker(d time.Duration) Ticker
	NewTimer(d time.Duration) Timer
}

var _ Clock = (*MockClock)(nil)
var _ Clock = (RealClock)(RealClock{})

type RealClock struct{}

func New() RealClock {
	return RealClock{}
}

func (RealClock) Now() time.Time {
	return time.Now()
}

func (RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (RealClock) AfterFunc(d time.Duration, f func()) Timer {
	return &RealTimer{Timer: time.AfterFunc(d, f)}
}

func (RealClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (RealClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (RealClock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

func (RealClock) NewTicker(d time.Duration) Ticker {
	return &RealTicker{Ticker: time.NewTicker(d)}
}

func (RealClock) NewTimer(d time.Duration) Timer {
	return &RealTimer{Timer: time.NewTimer(d)}
}

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
	return c.makeMockTicker(d, true, nil)
}

func (c *MockClock) NewTimer(d time.Duration) Timer {
	return c.makeMockTimer(d, false, nil)
}
