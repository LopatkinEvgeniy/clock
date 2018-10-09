# clock

[![Build Status](https://travis-ci.org/LopatkinEvgeniy/clock.png?branch=master)](https://travis-ci.org/LopatkinEvgeniy/clock) [![GoDoc](https://godoc.org/github.com/LopatkinEvgeniy/clock?status.svg)](http://godoc.org/github.com/LopatkinEvgeniy/clock)

Small timer-driven library for mocking time in Go. [Clockwork](https://github.com/jonboulle/clockwork) drop in replacement.

### Why another one time mocking library?
* Race free
* Full featured
* Redesigned

### Example

Suppose we have some type with time-dependent method that we wanna test.
Instead of direct use `time` package we specify the clock field:
```go
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
```

Now in tests we just inject FakeClock to the tested struct. 
This allows us to manipulate time:

```go
func TestExample(t *testing.T) {
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
```

In production simply inject the real clock instead
```go
mt := myType{clock: clock.NewRealClock()}
```

### Inspired by:
* https://github.com/jonboulle/clockwork
* https://github.com/benbjohnson/clock
* https://github.com/facebookgo/clock
