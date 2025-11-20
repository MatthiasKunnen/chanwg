package chanwg

import "sync"

// WaitGroup is a single-use synchronization primitive similar to [sync.WaitGroup].
//
// [sync.WaitGroup] exposes a blocking [sync.WaitGroup.Wait] method which cannot be abandoned
// without leaking a goroutine.
// [chanwg.WaitGroup], on the other hand, facilitates waiting using a channel which can be abandoned.
type WaitGroup struct {
	mu      sync.Mutex
	done    chan struct{}
	counter int
	state   state
}

type state uint8

const (
	stateInitial state = iota
	stateReadyForWait
	stateClosed
)

// Add adds delta, which may be negative, to the [WaitGroup] task counter.
// If the counter becomes zero, and Ready has been called, WaitChan closes.
// If the counter goes negative, Add panics.
//
// Callers should prefer [WaitGroup.Go].
//
// Add may happen before Ready or while the counter is not zero.
// Typically, Add is executed before the statement creating the goroutine or other event to be
// waited for.
// Typically, Add is executed before the statement creating the goroutine or other event to be
// waited for.
func (cwg *WaitGroup) Add(delta int) {
	if delta == 0 {
		return
	}

	cwg.mu.Lock()
	defer cwg.mu.Unlock()
	if cwg.state == stateClosed {
		panic("chanwg: WaitGroup already closed")
	}

	cwg.counter += delta

	switch {
	case cwg.counter < 0:
		panic("chanwg: negative WaitGroup counter, too many Done calls")
	case cwg.counter == 0:
		if cwg.state == stateReadyForWait {
			cwg.state = stateClosed
			if cwg.done != nil {
				close(cwg.done)
			}
		}
	}
}

// Done decrements the task counter by one.
// It is equivalent to Add(-1).
//
// Callers should prefer [WaitGroup.Go].
//
// When the counter reaches zero and Ready has been called, WaitChan is closed.
//
// Panics if:
//   - Done is called more times than Add
//   - WaitChan has closed
func (cwg *WaitGroup) Done() {
	cwg.Add(-1)
}

// WaitChan returns a channel that will be closed when all tasks have completed and Ready
// has been called.
// The channel can be received from before any tasks have been added.
// It will always return the same channel.
func (cwg *WaitGroup) WaitChan() <-chan struct{} {
	cwg.mu.Lock()
	defer cwg.mu.Unlock()
	if cwg.done == nil {
		cwg.done = make(chan struct{})
		if cwg.state == stateClosed {
			close(cwg.done)
		}
	}

	return cwg.done
}

// Ready causes the WaitChan to close when all tasks complete.
// It is typically called after all tasks are added.
// Ready may be called multiple times safely.
func (cwg *WaitGroup) Ready() {
	cwg.mu.Lock()
	defer cwg.mu.Unlock()
	if cwg.state != stateInitial {
		return
	}
	if cwg.counter == 0 {
		cwg.state = stateClosed
		if cwg.done != nil {
			close(cwg.done)
		}
	} else {
		cwg.state = stateReadyForWait
	}
}

// Go calls f in a new goroutine and adds that task to the WaitGroup.
// When f returns, the task is removed from the WaitGroup.
//
// The function f must not panic.
//
// Go may happen before Ready or while the WaitGroup is not empty.
// This means a goroutine started by Go may itself call Go.
func (cwg *WaitGroup) Go(f func()) {
	cwg.Add(1)
	go func() {
		defer cwg.Done()
		f()
	}()
}
