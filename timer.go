package clock

import (
	"time"
)

type Timer interface {
	Chan() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

var _ Timer = realTimer{}
var _ Timer = (*internalTimer)(nil)

type realTimer struct {
	*time.Timer
}

func (t realTimer) Chan() <-chan time.Time {
	return t.C
}
