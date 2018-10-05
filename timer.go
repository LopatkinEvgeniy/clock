package clock

import (
	"time"
)

type Timer interface {
	Ch() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

var _ Timer = (*RealTimer)(nil)
var _ Timer = (*internalTimer)(nil)

type RealTimer struct {
	*time.Timer
}

func (t *RealTimer) Ch() <-chan time.Time {
	return t.C
}
