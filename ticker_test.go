package clock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/LopatkinEvgeniy/clock"
)

func TestMockTickerCh(t *testing.T) {
	c := clock.NewMock()
	ticker := c.NewTicker(100 * time.Second)

	for i := 0; i < 100; i++ {
		for i := 0; i < 99; i++ {
			c.Add(time.Second)
			select {
			case <-ticker.Ch():
				t.Fatal("Unexpected ticker's channel receive")
			default:
			}
		}
		c.Add(time.Second)
		select {
		case <-ticker.Ch():
		default:
			t.Fatal("Expected receive from the ticker's channel")
		}
	}
}

func TestMockTickerChStress(t *testing.T) {
	c := clock.NewMock()
	ticker := c.NewTicker(10000 * time.Second)

	for i := 0; i < 100; i++ {
		wg := sync.WaitGroup{}
		for i := 0; i < 9999; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				c.Add(time.Second)
				select {
				case <-ticker.Ch():
					t.Fatal("Unexpected ticker's channel receive")
				default:
				}
			}()
		}
		wg.Wait()
		c.Add(time.Second)
		select {
		case <-ticker.Ch():
		default:
			t.Fatal("Expected receive from the ticker's channel")
		}
	}
}

func TestMockTickerStop(t *testing.T) {
	t.Run("stop before first tick", func(t *testing.T) {
		c := clock.NewMock()
		d := time.Nanosecond
		ticker := c.NewTicker(d)

		select {
		case <-ticker.Ch():
			t.Fatal("Unexpected ticker's channel receive")
		default:
		}

		wasActive := ticker.Stop()
		if !wasActive {
			t.Fatal("Unexpected stop result value")
		}

		for i := 0; i < 100; i++ {
			c.Add(d)
		}

		select {
		case <-ticker.Ch():
			t.Fatal("Unexpected ticker's channel receive")
		default:
		}
	})

	t.Run("stop after first tick", func(t *testing.T) {
		c := clock.NewMock()
		d := time.Nanosecond
		ticker := c.NewTicker(d)

		c.Add(d)

		select {
		case <-ticker.Ch():
		default:
			t.Fatal("Unexpected ticker's channel receive")
		}

		wasActive := ticker.Stop()
		if !wasActive {
			t.Fatal("Unexpected stop result value")
		}

		for i := 0; i < 100; i++ {
			c.Add(d)
		}

		select {
		case <-ticker.Ch():
			t.Fatal("Unexpected ticker's channel receive")
		default:
		}
	})

	t.Run("stop multiple times", func(t *testing.T) {
		c := clock.NewMock()
		ticker := c.NewTicker(time.Nanosecond)

		wasActive := ticker.Stop()
		if !wasActive {
			t.Fatal("Unexpected stop result value")
		}
		for i := 0; i < 5; i++ {
			wasActive := ticker.Stop()
			if wasActive {
				t.Fatal("Unexpected stop result value")
			}
		}

		c.Add(time.Hour)

		select {
		case <-ticker.Ch():
			t.Fatal("Unexpected ticker's channel receive")
		default:
		}
	})
}

func TestMockTickerStopStress(t *testing.T) {
	c := clock.NewMock()
	d := time.Hour

	for i := 0; i < 100000; i++ {
		ticker := c.NewTicker(d)

		go func() {
			c.Add(d)
		}()

		actualTime := <-ticker.Ch()
		expectedTime := c.Now()
		if expectedTime != actualTime {
			t.Fatalf("Unexpected time received from the channel, expected=%s, actual=%s", expectedTime, actualTime)
		}

		wasActive := ticker.Stop()
		if !wasActive {
			t.Fatal("Unexpected stop result value")
		}
	}
}
