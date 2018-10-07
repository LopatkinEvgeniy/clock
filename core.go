package clock

import (
	"sync"
	"time"
)

// internalTimer is an internal representaion
// of mock timers, tickers and sleepers.
type internalTimer struct {
	clock       *internalClock
	id          int
	ch          chan time.Time
	triggerTime time.Time
	isActive    bool
	callback    func()
	isTicker    bool
	duration    time.Duration
}

// internalClock in an internal implementation
// of base mock clock functionality.
// internalClock has it's own current time value.
// All active timers/tickers/waiters are registered here.
// All internalClock's methods are thread-safe.
type internalClock struct {
	mu          sync.Mutex
	now         time.Time
	nextTimerID int
	timers      map[int]*internalTimer
}

// newInternalClock creates a new initialized internalClock instance.
func newInternalClock(t time.Time) *internalClock {
	return &internalClock{
		now:    t,
		timers: make(map[int]*internalTimer),
	}
}

// getCurrentTime returns an internalClock's current time value.
func (c *internalClock) getCurrentTime() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

// moveTimeForward adds specified duration
// to the current internalClock's time.
// It will affect all registered timers.
func (c *internalClock) moveTimeForward(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = c.now.Add(d)

	for _, t := range c.timers {
		if t.triggerTime.After(c.now) {
			continue
		}

		if !t.isTicker {
			t.isActive = false
			delete(c.timers, t.id)
			if t.callback != nil {
				go t.callback()
			}
		}

		if t.callback == nil {
			select {
			case t.ch <- t.triggerTime:
			default:
			}
		}

		if t.isTicker {
			for !t.triggerTime.After(c.now) {
				t.triggerTime = t.triggerTime.Add(t.duration)
			}
		}
	}
}

// newInternalTimer creates and registres a new internalTimer instance.
func (c *internalClock) newInternalTimer(d time.Duration, isTicker bool, callback func()) *internalTimer {
	c.mu.Lock()
	defer c.mu.Unlock()

	if isTicker && callback != nil {
		panic("unexpected callback for the ticker")
	}

	t := &internalTimer{
		clock:       c,
		id:          c.nextTimerID,
		ch:          make(chan time.Time, 1),
		triggerTime: c.now.Add(d),
		isActive:    true,
		callback:    callback,
		isTicker:    isTicker,
		duration:    d,
	}
	c.timers[t.id] = t
	c.nextTimerID++

	return t
}

// stopTimer unregisters specified timer.
// It returns true if the specified timer was active.
func (c *internalClock) stopTimer(t *internalTimer) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.timers, t.id)

	timerWasActive := t.isActive
	t.isActive = false

	return timerWasActive
}

// resetTimer changes duration for the specified timer.
// Specified timer would be registered again.
// resetTimer returns true if specified timer was active.
func (c *internalClock) resetTimer(t *internalTimer, d time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if t.isTicker {
		panic("ticker cannot be reset")
	}

	timerWasActive := t.isActive
	t.isActive = true
	t.triggerTime = c.now.Add(d)
	t.duration = d

	if !timerWasActive {
		c.timers[t.id] = t
	}

	return timerWasActive
}

// waitersCount returns current count of registered timers, tickers and sleepers.
func (c *internalClock) waitersCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.timers)
}
