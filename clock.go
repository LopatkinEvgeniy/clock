package clock

import (
	"sync"
	"time"
)

type MockClock struct {
	mu          sync.Mutex
	now         time.Time
	nextTimerID int
	timers      map[int]*mockTimer
}

func NewMock() *MockClock {
	return &MockClock{
		timers: make(map[int]*mockTimer),
	}
}

func (c *MockClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

func (c *MockClock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = c.now.Add(d)

	for _, t := range c.timers {
		if t.triggerTime.After(c.now) {
			continue
		}

		select {
		case t.ch <- t.triggerTime:
		default:
		}

		if t.isTicker {
			for !t.triggerTime.After(c.now) {
				t.triggerTime = t.triggerTime.Add(t.duration)
			}
		} else {
			t.isActive = false
			delete(c.timers, t.id)
		}
	}
}

func (c *MockClock) NewTicker(d time.Duration) Ticker {
	return c.makeMockTimer(d, true)
}

func (c *MockClock) NewTimer(d time.Duration) Timer {
	return c.makeMockTimer(d, false)
}

func (c *MockClock) makeMockTimer(d time.Duration, isTicker bool) *mockTimer {
	c.mu.Lock()
	defer c.mu.Unlock()

	t := &mockTimer{
		clock:       c,
		id:          c.nextTimerID,
		ch:          make(chan time.Time, 1),
		triggerTime: c.now.Add(d),
		isActive:    true,

		isTicker: isTicker,
		duration: d,
	}
	c.timers[t.id] = t
	c.nextTimerID++

	return t
}

func (c *MockClock) stopTimer(t *mockTimer) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.timers, t.id)

	timerWasActive := t.isActive
	t.isActive = false

	return timerWasActive
}

func (c *MockClock) resetTimer(t *mockTimer, d time.Duration) bool {
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
