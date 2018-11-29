package clock

import (
	"time"
)

// Clock contains various time's functions that can be mocked.
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

var _ Clock = FakeClock{}
var _ Clock = realClock{}

// realClock is just a time's wrapper.
type realClock struct{}

// NewRealClock returns a new instance of the real clock.
func NewRealClock() Clock {
	return realClock{}
}

// Now implements Clock.
func (realClock) Now() time.Time {
	return time.Now()
}

// After implements Clock.
func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// AfterFunc implements Clock.
func (realClock) AfterFunc(d time.Duration, f func()) Timer {
	return realTimer{Timer: time.AfterFunc(d, f)}
}

// Since implements Clock.
func (realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Until implements Clock.
func (realClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

// Sleep implements Clock.
func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Tick implements Clock.
func (realClock) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

// NewTicker implements Clock.
// It returns a new instance of the real ticker.
func (realClock) NewTicker(d time.Duration) Ticker {
	return realTicker{Ticker: time.NewTicker(d)}
}

// NewTimer implements Clock.
// It returns a new instance of the real timer.
func (realClock) NewTimer(d time.Duration) Timer {
	return realTimer{Timer: time.NewTimer(d)}
}

// FakeClock is an internalClock's shallow wrapper.
// It provides special mock methods such Advance or WaitersCount.
type FakeClock struct {
	*internalClock
}

// NewFakeClock returns a new instance of the fake clock.
func NewFakeClock() FakeClock {
	return FakeClock{
		internalClock: newInternalClock(time.Time{}),
	}
}

// NewFakeClockAt returns a new instance of the fake clock.
// Specified time will be used as a current clock's time.
func NewFakeClockAt(t time.Time) FakeClock {
	return FakeClock{
		internalClock: newInternalClock(t),
	}
}

// Advance moves current clock's time forward.
// It affects all active timers/tickers/sleepers.
func (c FakeClock) Advance(d time.Duration) {
	c.moveTimeForward(d)
}

// WaitersCount returns current active timers/tickers/sleepers count.
func (c FakeClock) WaitersCount() int {
	return c.waitersCount()
}

// BlockUntil waits for the specified count of active timers/tickers/sleepers.
func (c FakeClock) BlockUntil(n int) {
	for c.waitersCount() != n {
		time.Sleep(time.Millisecond)
	}
}

// Now implements Clock.
func (c FakeClock) Now() time.Time {
	return c.getCurrentTime()
}

// After implements Clock.
func (c FakeClock) After(d time.Duration) <-chan time.Time {
	return c.NewTimer(d).Chan()
}

// AfterFunc implements Clock.
func (c FakeClock) AfterFunc(d time.Duration, f func()) Timer {
	return mockTimer{
		internalTimer: c.newInternalTimer(d, false, f),
	}
}

// Since implements Clock.
func (c FakeClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

// Until implements Clock.
func (c FakeClock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}

// Sleep implements Clock.
func (c FakeClock) Sleep(d time.Duration) {
	<-c.NewTimer(d).Chan()
}

// Tick implements Clock.
func (c FakeClock) Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return c.NewTicker(d).Chan()
}

// NewTicker implements Clock.
// It returns a new instance of the mock ticker.
func (c FakeClock) NewTicker(d time.Duration) Ticker {
	if d <= 0 {
		panic("non-positive interval for NewTicker")
	}
	return mockTicker{
		internalTimer: c.newInternalTimer(d, true, nil),
	}
}

// NewTimer implements Clock.
// It returns a new instance of the mock timer.
func (c FakeClock) NewTimer(d time.Duration) Timer {
	return mockTimer{
		internalTimer: c.newInternalTimer(d, false, nil),
	}
}
