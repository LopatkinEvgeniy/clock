package clock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/LopatkinEvgeniy/clock"
)

func TestFakeTickerChan(t *testing.T) {
	c := clock.NewFakeClock()
	ticker := c.NewTicker(100 * time.Second)

	for i := 0; i < 100; i++ {
		for i := 0; i < 99; i++ {
			c.Advance(time.Second)
			select {
			case <-ticker.Chan():
				t.Fatal("Unexpected ticker's channel receive")
			default:
			}
		}
		c.Advance(time.Second)
		select {
		case <-ticker.Chan():
		default:
			t.Fatal("Expected receive from the ticker's channel")
		}
	}
}

func TestFakeTickerChanStress(t *testing.T) {
	c := clock.NewFakeClock()
	ticker := c.NewTicker(10000 * time.Second)

	for i := 0; i < 100; i++ {
		wg := sync.WaitGroup{}
		for i := 0; i < 9999; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				c.Advance(time.Second)
				select {
				case <-ticker.Chan():
					t.Fatal("Unexpected ticker's channel receive")
				default:
				}
			}()
		}
		wg.Wait()
		c.Advance(time.Second)
		select {
		case <-ticker.Chan():
		default:
			t.Fatal("Expected receive from the ticker's channel")
		}
	}
}

func TestFakeTickerStop(t *testing.T) {
	t.Run("stop before first tick", func(t *testing.T) {
		c := clock.NewFakeClock()
		d := time.Nanosecond
		ticker := c.NewTicker(d)

		select {
		case <-ticker.Chan():
			t.Fatal("Unexpected ticker's channel receive")
		default:
		}

		ticker.Stop()

		for i := 0; i < 100; i++ {
			c.Advance(d)
		}

		select {
		case <-ticker.Chan():
			t.Fatal("Unexpected ticker's channel receive")
		default:
		}
	})

	t.Run("stop after first tick", func(t *testing.T) {
		c := clock.NewFakeClock()
		d := time.Nanosecond
		ticker := c.NewTicker(d)

		c.Advance(d)

		select {
		case <-ticker.Chan():
		default:
			t.Fatal("Unexpected ticker's channel receive")
		}

		ticker.Stop()

		for i := 0; i < 100; i++ {
			c.Advance(d)
		}

		select {
		case <-ticker.Chan():
			t.Fatal("Unexpected ticker's channel receive")
		default:
		}
	})

	t.Run("stop multiple times", func(t *testing.T) {
		c := clock.NewFakeClock()
		ticker := c.NewTicker(time.Nanosecond)

		for i := 0; i < 5; i++ {
			ticker.Stop()
		}

		c.Advance(time.Hour)

		select {
		case <-ticker.Chan():
			t.Fatal("Unexpected ticker's channel receive")
		default:
		}
	})
}

func TestFakeTickerStopStress(t *testing.T) {
	c := clock.NewFakeClock()
	d := time.Hour

	for i := 0; i < 100000; i++ {
		ticker := c.NewTicker(d)

		go func() {
			c.Advance(d)
		}()

		actualTime := <-ticker.Chan()
		expectedTime := c.Now()
		if expectedTime != actualTime {
			t.Fatalf("Unexpected time received from the channel, expected=%s, actual=%s", expectedTime, actualTime)
		}

		ticker.Stop()
	}
}

func TestFakeTickerPanicsOnNonPositiveInterval(t *testing.T) {
	c := clock.NewFakeClock()

	t.Run("zero interval", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected a panic")
			}
		}()
		c.NewTicker(0)
	})

	t.Run("negative interval", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected a panic")
			}
		}()
		c.NewTicker(-1)
	})
}
