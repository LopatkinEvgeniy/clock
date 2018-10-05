package clock

import "time"

type Ticker interface {
	Ch() <-chan time.Time
	Stop()
}

var _ Ticker = (*RealTicker)(nil)
var _ Ticker = (*internalTicker)(nil)

type RealTicker struct {
	*time.Ticker
}

func (t *RealTicker) Ch() <-chan time.Time {
	return t.C
}
