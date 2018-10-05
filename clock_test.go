package clock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/LopatkinEvgeniy/clock"
)

func TestFakeClockNow(t *testing.T) {
	t.Run("use NewFakeClock constructor", func(t *testing.T) {
		c := clock.NewFakeClock()
		expected := time.Time{}

		for i := 0; i < 100; i++ {
			c.Add(time.Hour)
			expected = expected.Add(time.Hour)

			now := c.Now()
			if now != expected {
				t.Fatalf("unexpected now result, expected: %s, actual: %s", expected, now)
			}
		}
	})

	t.Run("use NewFakeClockAt constructor", func(t *testing.T) {
		initialTime := time.Now()
		c := clock.NewFakeClockAt(initialTime)
		expected := initialTime

		for i := 0; i < 100; i++ {
			c.Add(time.Hour)
			expected = expected.Add(time.Hour)

			now := c.Now()
			if now != expected {
				t.Fatalf("unexpected now result, expected: %s, actual: %s", expected, now)
			}
		}
	})
}

func TestFakeClockNowStress(t *testing.T) {
	c := clock.NewFakeClock()

	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			c.Add(time.Minute)
			wg.Done()
		}()
	}
	wg.Wait()

	expected := (time.Time{}).Add(100000 * time.Minute)
	now := c.Now()
	if now != expected {
		t.Fatalf("unexpected now result, expected: %s, actual: %s", expected, now)
	}
}

func TestFakeClockAfter(t *testing.T) {
	c := clock.NewFakeClock()
	ch := c.After(100 * time.Minute)

	for i := 0; i < 99; i++ {
		c.Add(time.Minute)

		select {
		case <-ch:
			t.Fatal("Unexpected channel receive")
		default:
		}
	}

	c.Add(time.Minute)
	var actualTime time.Time

	select {
	case actualTime = <-ch:
	default:
		t.Fatal("Expected channel receive")
	}

	expectedTime := (time.Time{}).Add(100 * time.Minute)
	if expectedTime != actualTime {
		t.Fatalf("Unexpected time received from the channel, expected=%s, actual=%s", expectedTime, actualTime)
	}
}

// TODO: AfterFunc tests
