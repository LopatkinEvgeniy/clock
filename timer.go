package clock

import (
	"time"
)

type Timer interface {
	Ch() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}
