package clock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/LopatkinEvgeniy/clock"
)

func TestFakeTimerCh(t *testing.T) {
	c := clock.NewFakeClock()
	timer := c.NewTimer(100 * time.Second)

	for i := 0; i < 99; i++ {
		c.Add(time.Second)
		select {
		case <-timer.Chan():
			t.Fatal("Unexpected timer's channel receive")
		default:
		}
	}

	c.Add(time.Second)
	select {
	case <-timer.Chan():
	default:
		t.Fatal("Expected receive from the timer's channel")
	}
}

func TestFakeTimerChStress(t *testing.T) {
	c := clock.NewFakeClock()
	timer := c.NewTimer(10000 * time.Second)

	wg := sync.WaitGroup{}
	for i := 0; i < 9999; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			c.Add(time.Second)
			select {
			case <-timer.Chan():
				t.Fatal("Unexpected timer's channel receive")
			default:
			}
		}()
	}
	wg.Wait()

	c.Add(time.Second)
	select {
	case <-timer.Chan():
	default:
		t.Fatal("Expected receive from the timer's channel")
	}
}

func TestFakeTimerReset(t *testing.T) {
	t.Run("timer not expired", func(t *testing.T) {
		c := clock.NewFakeClock()
		timer := c.NewTimer(time.Nanosecond)

		select {
		case <-timer.Chan():
			t.Fatal("Unexpected timer's channel receive")
		default:
		}

		wasActive := timer.Reset(time.Second)
		if !wasActive {
			t.Fatal("Unexpected reset result value")
		}
	})

	t.Run("timer expired", func(t *testing.T) {
		c := clock.NewFakeClock()
		d := time.Minute
		timer := c.NewTimer(d)

		c.Add(d)

		select {
		case <-timer.Chan():
		default:
			t.Fatal("Expected receive from the timer's channel")
		}

		wasActive := timer.Reset(time.Second)
		if wasActive {
			t.Fatal("Unexpected reset result value")
		}
	})

	t.Run("reset multiple times", func(t *testing.T) {
		c := clock.NewFakeClock()
		timer := c.NewTimer(time.Nanosecond)
		expectedDuration := time.Minute

		timer.Reset(10 * time.Minute)
		timer.Reset(5 * time.Minute)
		timer.Reset(expectedDuration)

		c.Add(expectedDuration)

		actualTime := <-timer.Chan()
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

func TestFakeTimerResetStress(t *testing.T) {
	c := clock.NewFakeClock()
	d := time.Hour
	timer := c.NewTimer(d)

	for i := 0; i < 100000; i++ {
		go func() {
			c.Add(d)
		}()

		actualTime := <-timer.Chan()
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

func TestFakeTimerStop(t *testing.T) {
	t.Run("timer not expired", func(t *testing.T) {
		c := clock.NewFakeClock()
		timer := c.NewTimer(time.Nanosecond)

		select {
		case <-timer.Chan():
			t.Fatal("Unexpected timer's channel receive")
		default:
		}

		wasActive := timer.Stop()
		if !wasActive {
			t.Fatal("Unexpected stop result value")
		}

		c.Add(time.Hour)

		select {
		case <-timer.Chan():
			t.Fatal("Unexpected timer's channel receive")
		default:
		}
	})

	t.Run("timer expired", func(t *testing.T) {
		c := clock.NewFakeClock()
		d := time.Minute
		timer := c.NewTimer(d)

		c.Add(d)

		select {
		case <-timer.Chan():
		default:
			t.Fatal("Expected receive from the timer's channel")
		}

		wasActive := timer.Stop()
		if wasActive {
			t.Fatal("Unexpected stop result value")
		}

		select {
		case <-timer.Chan():
			t.Fatal("Unexpected timer's channel receive")
		default:
		}
	})

	t.Run("stop multiple times", func(t *testing.T) {
		c := clock.NewFakeClock()
		timer := c.NewTimer(time.Nanosecond)

		wasActive := timer.Stop()
		if !wasActive {
			t.Fatal("Unexpected stop result value")
		}
		for i := 0; i < 5; i++ {
			wasActive := timer.Stop()
			if wasActive {
				t.Fatal("Unexpected stop result value")
			}
		}

		c.Add(time.Hour)

		select {
		case <-timer.Chan():
			t.Fatal("Unexpected timer's channel receive")
		default:
		}
	})
}

func TestFakeTimerStopStress(t *testing.T) {
	c := clock.NewFakeClock()
	d := time.Hour

	for i := 0; i < 100000; i++ {
		timer := c.NewTimer(d)

		go func() {
			c.Add(d)
		}()

		actualTime := <-timer.Chan()
		expectedTime := c.Now()
		if expectedTime != actualTime {
			t.Fatalf("Unexpected time received from the channel, expected=%s, actual=%s", expectedTime, actualTime)
		}

		wasActive := timer.Stop()
		if wasActive {
			t.Fatal("Unexpected stop result value")
		}
	}
}
