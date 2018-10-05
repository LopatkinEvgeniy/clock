package clock

import "time"

type Ticker interface {
	Ch() <-chan time.Time
	Stop() bool
}
