package clock

import (
	"sync"
	"time"
)

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

type internalClock struct {
	mu          sync.Mutex
	now         time.Time
	nextTimerID int
	timers      map[int]*internalTimer
}

func newInternalClock(t time.Time) *internalClock {
	return &internalClock{
		now:    t,
		timers: make(map[int]*internalTimer),
	}
}

func (c *internalClock) getCurrentTime() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

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

		isTicker: isTicker,
		duration: d,
	}
	c.timers[t.id] = t
	c.nextTimerID++

	return t
}

func (c *internalClock) stopTimer(t *internalTimer) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.timers, t.id)

	timerWasActive := t.isActive
	t.isActive = false

	return timerWasActive
}

func (c *internalClock) resetTimer(t *internalTimer, d time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if t.isTicker {
		panic("ticker cannot be reset")
	}

	timerWasActive := t.isActive
	t.isActive = true
	t.triggerTime = c.now.Add(d)

	if !timerWasActive {
		c.timers[t.id] = t
	}

	return timerWasActive
}

func (c *internalClock) waitersCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.timers)
}
