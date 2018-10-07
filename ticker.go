package clock

import "time"

// Ticker is an interface that represents both real and mock tickers.
type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}

var _ Ticker = realTicker{}
var _ Ticker = mockTicker{}

// realTicker is just a time.Ticker's shallow wrapper.
type realTicker struct {
	*time.Ticker
}

// Chan implementats Ticker.
func (t realTicker) Chan() <-chan time.Time {
	return t.C
}

// mockTicker is just an internalTimer's shallow wrapper.
type mockTicker struct {
	*internalTimer
}

// Chan implementats Ticker.
func (t mockTicker) Chan() <-chan time.Time {
	return t.ch
}

// Stop implementats Ticker.
func (t mockTicker) Stop() {
	t.clock.stopTimer(t.internalTimer)
}
