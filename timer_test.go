package clock_test

import (
	"testing"
	"time"

	"github.com/LopatkinEvgeniy/clock"
)

func TestMockTimerReset(t *testing.T) {
	t.Run("timer not expired", func(t *testing.T) {
		c := clock.NewMock()
		timer := c.NewTimer(time.Nanosecond)

		select {
		case <-timer.Ch():
			t.Fatal("Unexpected timer's channel receive")
		default:
		}

		wasActive := timer.Reset(time.Second)
		if !wasActive {
			t.Fatal("Unexpected reset result value")
		}
	})

	t.Run("timer expired", func(t *testing.T) {
		c := clock.NewMock()
		d := time.Minute
		timer := c.NewTimer(d)

		c.Add(d)

		select {
		case <-timer.Ch():
		default:
			t.Fatal("Expected receive from the timer's channel")
		}

		wasActive := timer.Reset(time.Second)
		if wasActive {
			t.Fatal("Unexpected reset result value")
		}
	})

	t.Run("reset multiple times", func(t *testing.T) {
		c := clock.NewMock()
		timer := c.NewTimer(time.Nanosecond)
		expectedDuration := time.Minute

		timer.Reset(10 * time.Minute)
		timer.Reset(5 * time.Minute)
		timer.Reset(expectedDuration)

		c.Add(expectedDuration)

		actualTime := <-timer.Ch()
		expectedTime := c.Now()
		if expectedTime != actualTime {
			t.Fatalf("Unexpected time received from the channel, expected=%s, actual=%s", expectedTime, actualTime)
		}

		wasActive := timer.Reset(time.Hour)
		if wasActive {
			t.Fatal("Unexpected reset result value")
		}
	})
}

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
