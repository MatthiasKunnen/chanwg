package chanwg

import "sync"

// WaitGroup is a single-use synchronization primitive similar to [sync.WaitGroup].
//
// Instead of a blocking Wait method, it exposes a channel that closes when all tracked
// operations have completed.
//
// Waiting for a [sync.WaitGroup] that never completes blocks the goroutine indefinitely.
// [chanwg.WaitGroup], on the other hand, allows abandoning a wait.
//
// WaitGroup requires at least one call to Go (or Add and corresponding Done) before completing.
// This allows you to extract the channel before calling Add or Go.
type WaitGroup struct {
	counter int
	closed  bool
	mu      sync.Mutex
	done    chan struct{}
}

// Add increments the counter by the given positive delta.
// Add must be called before the corresponding operations begin execution.
//
// Callers should prefer [WaitGroup.Go].
//
// Add may happen while the task counter has never reached zero after the initial state.
// Typically, Add is executed before the statement creating the goroutine or other event to be
// waited for.
func (cwg *WaitGroup) Add(delta int) {
	if delta == 0 {
		return
	}

	cwg.mu.Lock()
	defer cwg.mu.Unlock()
	if cwg.closed {
		panic("chanwg: WaitGroup already closed")
	}

	cwg.counter += delta

	switch {
	case cwg.counter < 0:
		panic("chanwg: negative WaitGroup counter, too many Done calls")
	case cwg.counter == 0:
		cwg.closed = true
		if cwg.done != nil {
			close(cwg.done)
		}
	}
}

// Done decrements the task counter by one.
// It is equivalent to Add(-1).
//
// Callers should prefer [WaitGroup.Go].
//
// When the counter reaches zero has been called, WaitChan is closed.
//
// Panics if:
//   - Done is called more times than Add
//   - WaitChan has closed
func (cwg *WaitGroup) Done() {
	cwg.Add(-1)
}

// WaitChan returns a channel that will be closed when all tracked operations are complete.
func (cwg *WaitGroup) WaitChan() <-chan struct{} {
	cwg.mu.Lock()
	defer cwg.mu.Unlock()
	if cwg.done == nil {
		cwg.done = make(chan struct{})
		if cwg.closed {
			close(cwg.done)
		}
	}

	return cwg.done
}

// Go calls f in a new goroutine and adds that task to the WaitGroup.
// When f returns, the task is removed from the WaitGroup.
//
// The function f must not panic.
//
// Go may happen while the WaitGroup has not become empty.
// This means a goroutine started by Go may itself call Go.
func (cwg *WaitGroup) Go(f func()) {
	cwg.Add(1)
	go func() {
		defer cwg.Done()
		f()
	}()
}
