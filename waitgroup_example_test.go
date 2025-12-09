package chanwg_test

import (
	"context"
	"fmt"
	"time"

	"github.com/MatthiasKunnen/chanwg/v2"
)

// This example shows how to await the WaitGroup and another channel simultaneously.
func ExampleWaitGroup_basic() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()

	var wg chanwg.WaitGroup
	wg.Go(func() {
		select {} // Wait forever
	})
	wg.Ready()

	select {
	case <-wg.WaitChan():
	case <-ctx.Done():
		fmt.Println("Abort start")
	}

	// Output: Abort start
}

// The WaitChan can be received from even before adding tasks.
// It won't complete before Ready and the tasks being done.
func ExampleWaitGroup_useWaitChanBeforeAddingTasks() {
	var wg chanwg.WaitGroup

	go func() {
		// Time this function after the WaitChan receive
		time.Sleep(10 * time.Millisecond)
		wg.Go(func() {
			// Do work
		})
		wg.Ready()
	}()

	select {
	case <-wg.WaitChan():
		fmt.Println("Done")
	case <-time.After(50 * time.Millisecond):
		fmt.Println("WaitChan did not complete as expected")
	}

	// Output: Done
}
