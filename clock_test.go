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
			c.Advance(time.Hour)
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
			c.Advance(time.Hour)
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
			c.Advance(time.Minute)
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
		c.Advance(time.Minute)

		select {
		case <-ch:
			t.Fatal("Unexpected channel receive")
		default:
		}
	}

	c.Advance(time.Minute)
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

func TestFakeClockAfterFunc(t *testing.T) {
	t.Run("check callback", func(t *testing.T) {
		c := clock.NewFakeClock()

		doneCh := make(chan struct{})
		timer := c.AfterFunc(10*time.Minute, func() {
			close(doneCh)
		})

		c.Advance(5 * time.Minute)
		time.Sleep(100 * time.Millisecond)

		select {
		case <-doneCh:
			t.Fatal("Permature callback's call")
		default:
		}

		c.Advance(5 * time.Minute)
		<-doneCh

		time.Sleep(100 * time.Millisecond)
		select {
		case <-timer.Chan():
			t.Fatal("Unexpected receive from the timer's channel")
		default:
		}
	})

	t.Run("stop timer", func(t *testing.T) {
		c := clock.NewFakeClock()

		doneCh := make(chan struct{})
		timer := c.AfterFunc(10*time.Minute, func() {
			close(doneCh)
		})

		wasActive := timer.Stop()
		if !wasActive {
			t.Fatal("Unexpected stop result value")
		}

		for i := 0; i < 10; i++ {
			c.Advance(10 * time.Minute)
			time.Sleep(10 * time.Millisecond)
			select {
			case <-doneCh:
				t.Fatal("Unexpected callback after Stop call")
			default:
			}
		}
	})
}

func TestFakeClockSince(t *testing.T) {
	tm := (time.Time{}).Add(time.Minute)

	c := clock.NewFakeClock()
	c.Advance(10 * time.Minute)

	expectedSince := 9 * time.Minute
	actualSince := c.Since(tm)
	if expectedSince != actualSince {
		t.Fatalf("Unexpected Since result value, expected=%s, actual=%s", expectedSince, actualSince)
	}

	c.Advance(10 * time.Second)
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

		c.Advance(time.Minute)
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

	c.BlockUntil(3)

	select {
	case <-ch1:
		t.Fatal("sleep ends too soon")
	case <-ch2:
		t.Fatal("sleep ends too soon")
	case <-ch3:
		t.Fatal("sleep ends too soon")
	default:
	}

	c.Advance(time.Minute)

	select {
	case <-ch1:
	case <-ch2:
		t.Fatal("sleep ends too soon")
	case <-ch3:
		t.Fatal("sleep ends too soon")
	}

	c.Advance(time.Minute)

	select {
	case <-ch2:
	case <-ch3:
		t.Fatal("sleep ends too soon")
	}

	c.Advance(time.Minute)
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
			c.Advance(time.Second)
			select {
			case <-tickCh:
				t.Fatal("Unexpected ticker's channel receive")
			default:
			}
		}
		c.Advance(time.Second)
		select {
		case <-tickCh:
		default:
			t.Fatal("Expected receive from the ticker's channel")
		}
	}
}
