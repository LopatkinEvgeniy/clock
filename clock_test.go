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

func TestFakeClockSince(t *testing.T) {
	tm := (time.Time{}).Add(time.Minute)

	c := clock.NewFakeClock()
	c.Add(10 * time.Minute)

	expectedSince := 9 * time.Minute
	actualSince := c.Since(tm)
	if expectedSince != actualSince {
		t.Fatalf("Unexpected Since result value, expected=%s, actual=%s", expectedSince, actualSince)
	}

	c.Add(10 * time.Second)
	expectedSince = expectedSince + 10*time.Second
	actualSince = c.Since(tm)
	if expectedSince != actualSince {
		t.Fatalf("Unexpected Since result value, expected=%s, actual=%s", expectedSince, actualSince)
	}
}

func TestFakeClockUntil(t *testing.T) {
	tm := (time.Time{}).Add(10 * time.Minute)

	c := clock.NewFakeClock()

	expectedUntil := 10 * time.Minute

	for i := 0; i < 9; i++ {
		actualUntil := c.Until(tm)
		if expectedUntil != actualUntil {
			t.Fatalf("Unexpected Until result value, expected=%s, actual=%s", expectedUntil, actualUntil)
		}

		c.Add(time.Minute)
		expectedUntil -= time.Minute
	}
}

func TestFakeClockSleep(t *testing.T) {
	c := clock.NewFakeClock()

	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})

	go func() {
		c.Sleep(time.Minute)
		close(ch1)
	}()
	go func() {
		c.Sleep(2 * time.Minute)
		close(ch2)
	}()
	go func() {
		c.Sleep(3 * time.Minute)
		close(ch3)
	}()

	for {
		if c.WaitersCount() == 3 {
			break
		}
		time.Sleep(time.Millisecond)
	}

	select {
	case <-ch1:
		t.Fatal("sleep ends too soon")
	case <-ch2:
		t.Fatal("sleep ends too soon")
	case <-ch3:
		t.Fatal("sleep ends too soon")
	default:
	}

	c.Add(time.Minute)

	select {
	case <-ch1:
	case <-ch2:
		t.Fatal("sleep ends too soon")
	case <-ch3:
		t.Fatal("sleep ends too soon")
	}

	c.Add(time.Minute)

	select {
	case <-ch2:
	case <-ch3:
		t.Fatal("sleep ends too soon")
	}

	c.Add(time.Minute)
	<-ch3
}

func TestFakeClockTick(t *testing.T) {
	c := clock.NewFakeClock()

	if c.Tick(0) != nil {
		t.Fatal("Nil channel expected")
	}
	if c.Tick(-1*time.Minute) != nil {
		t.Fatal("Nil channel expected")
	}

	tickCh := c.Tick(100 * time.Second)

	for i := 0; i < 100; i++ {
		for i := 0; i < 99; i++ {
			c.Add(time.Second)
			select {
			case <-tickCh:
				t.Fatal("Unexpected ticker's channel receive")
			default:
			}
		}
		c.Add(time.Second)
		select {
		case <-tickCh:
		default:
			t.Fatal("Expected receive from the ticker's channel")
		}
	}
}
