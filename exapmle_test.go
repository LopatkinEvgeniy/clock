package clock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/LopatkinEvgeniy/clock"
)

const incrementStateDelay = time.Hour

// myType is a type with time-dependent method that we will test.
type myType struct {
	clock clock.Clock
	state int
}

// incrementState increments myType's state with delay.
func (f *myType) incrementState() {
	f.clock.Sleep(incrementStateDelay)
	f.state++
}

func TestFeature(t *testing.T) {
	fakeClock := clock.NewFakeClock()

	// create the myType instance with fake clock.
	mt := myType{clock: fakeClock}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		mt.incrementState()
		wg.Done()
	}()

	// wait until incrementState is actually sleeping.
	fakeClock.BlockUntil(1)

	// assert state not changed.
	if mt.state != 0 {
		t.Fatalf("Unxepected state, expected=0 actual=%d", mt.state)
	}

	// move time forward and wait for incrementState done.
	fakeClock.Advance(incrementStateDelay)
	wg.Wait()

	// assert state incremented.
	if mt.state != 1 {
		t.Fatalf("Unxepected state, expected=1 actual=%d", mt.state)
	}
}
