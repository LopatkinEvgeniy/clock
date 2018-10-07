package clock

import "time"

type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}

var _ Ticker = realTicker{}
var _ Ticker = mockTicker{}

type realTicker struct {
	*time.Ticker
}

func (t realTicker) Chan() <-chan time.Time {
	return t.C
}

type mockTicker struct {
	*internalTimer
}

func (t mockTicker) Chan() <-chan time.Time {
	return t.ch
}

func (t mockTicker) Stop() {
	t.clock.stopTimer(t.internalTimer)
}
