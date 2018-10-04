package clock_test

import (
	"testing"
	"time"

	"github.com/LopatkinEvgeniy/clock"
)

func TestMockTimerResetStress(t *testing.T) {
	c := clock.NewMock()
	d := time.Hour
	timer := c.NewTimer(d)

	for i := 0; i < 100000; i++ {
		go func() {
			c.Add(d)
		}()

		actualTime := <-timer.Ch()
		expectedTime := c.Now()
		if expectedTime != actualTime {
			t.Fatalf("Unexpected time received from the channel, expected=%s, actual=%s", expectedTime, actualTime)
		}

		wasActive := timer.Reset(d)
		if wasActive {
			t.Fatal("Unexpected reset result value")
		}
	}
}
