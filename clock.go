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

var _ Clock = (FakeClock)(FakeClock{})
var _ Clock = (realClock)(realClock{})

type realClock struct{}

func NewRealClock() Clock {
	return realClock{}
}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (realClock) AfterFunc(d time.Duration, f func()) Timer {
	return &realTimer{Timer: time.AfterFunc(d, f)}
}

func (realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (realClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (realClock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

func (realClock) NewTicker(d time.Duration) Ticker {
	return &realTicker{Ticker: time.NewTicker(d)}
}

func (realClock) NewTimer(d time.Duration) Timer {
	return &realTimer{Timer: time.NewTimer(d)}
}

type FakeClock struct {
	*internalClock
}

func NewFakeClock() FakeClock {
	return FakeClock{
		internalClock: newInternalClock(time.Time{}),
	}
}

func NewFakeClockAt(t time.Time) FakeClock {
	return FakeClock{
		internalClock: newInternalClock(t),
	}
}

func (c FakeClock) Add(d time.Duration) {
	c.moveTimeForward(d)
}

func (c FakeClock) Now() time.Time {
	return c.getCurrentTime()
}

func (c FakeClock) After(d time.Duration) <-chan time.Time {
	return c.NewTimer(d).Chan()
}

func (c FakeClock) AfterFunc(d time.Duration, f func()) Timer {
	return c.newInternalTimer(d, false, f)
}

func (c FakeClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

func (c FakeClock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}

func (c FakeClock) Sleep(d time.Duration) {
	<-c.NewTimer(d).Chan()
}

func (c FakeClock) Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return c.NewTicker(d).Chan()
}

func (c FakeClock) NewTicker(d time.Duration) Ticker {
	return c.newInternalTicker(d, true, nil)
}

func (c FakeClock) NewTimer(d time.Duration) Timer {
	return c.newInternalTimer(d, false, nil)
}
