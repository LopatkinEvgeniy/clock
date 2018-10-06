package clock

import "time"

type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}

var _ Ticker = realTicker{}
var _ Ticker = (*internalTicker)(nil)

type realTicker struct {
	*time.Ticker
}

func (t realTicker) Chan() <-chan time.Time {
	return t.C
}
