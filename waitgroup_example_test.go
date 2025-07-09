package chanwg_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/MatthiasKunnen/chanwg"
	"time"
)

func ExampleWaitGroup_basic() {
	var wg chanwg.WaitGroup
	wg.Add(1)

	select {
	case <-wg.WaitChan():
	case <-time.After(time.Millisecond):
		fmt.Println("Abort start")
	}

	// Output: Abort start
}

type Foo struct {
	startedWg chanwg.WaitGroup
	started   <-chan struct{}
}

func NewFoo() *Foo {
	foo := &Foo{}
	foo.started = foo.startedWg.WaitChan()

	return foo
}

func (f *Foo) Start(ctx context.Context) error {
	f.startedWg.Add(1)

	go func() {
		time.Sleep(200 * time.Millisecond)
		defer f.startedWg.Done()
	}()

	select {
	case <-f.started:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (f *Foo) IsStarted() bool {
	select {
	case <-f.started:
		return true
	default:
		return false
	}
}

// Showcases how to use a [context.Context] to cancel a [chanwg.WaitGroup].
// In this example, the WaitGroup completes before the deadline.
func ExampleWaitGroup_struct_context_in_time() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFunc()
	err := NewFoo().Start(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("Failed to start, too slow")
		return
	} else if err != nil {
		fmt.Printf("Failed to start: %s\n", err)
		return
	}

	fmt.Println("Started in time")

	// Output: Started in time
}

// Showcases how to use a [context.Context] to cancel a [chanwg.WaitGroup].
// In this example, the context is canceled before the WaitGroup completes.
func ExampleWaitGroup_struct_context_canceled() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()
	err := NewFoo().Start(ctx)
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("Failed to start, too slow")
		return
	} else if err != nil {
		fmt.Printf("Failed to start: %s\n", err)
		return
	}

	fmt.Println("Started in time")

	// Output: Failed to start, too slow
}
