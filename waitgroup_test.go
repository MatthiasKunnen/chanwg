package chanwg_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MatthiasKunnen/chanwg/v2"
)

const tooManyDoneCallsPanic = "chanwg: negative WaitGroup counter, too many Done calls"
const alreadyClosedPanic = "chanwg: WaitGroup already closed"

func TestWaitGroupBasic(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	wg.Add(1)

	select {
	case <-wg.WaitChan():
		t.Fatal("WaitChan should not be closed yet")
	case <-time.After(100 * time.Millisecond):
		// Expected, give goroutine some time to potentially complete if there was a bug
	}

	wg.Done()

	select {
	case <-wg.WaitChan():
		t.Fatal("WaitChan should not be closed before Ready")
	case <-time.After(100 * time.Millisecond):
		// Expected, give goroutine some time to potentially complete if there was a bug
	}

	wg.Ready()

	select {
	case <-wg.WaitChan():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitChan was not closed after Done was called")
	}
}

func TestWaitGroupReadyBeforeDone(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	wg.Add(3)

	go func() {
		wg.Done()
	}()
	go func() {
		wg.Done()
	}()

	wg.Ready()

	select {
	case <-wg.WaitChan():
		t.Fatal("WaitChan should not be closed yet with 1 pending")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}

	wg.Done() // This should make the counter zero

	select {
	case <-wg.WaitChan():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitChan was not closed after all Dones were called")
	}
}

func TestWaitGroupReadyAfterDone(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	wg.Add(3)

	go func() {
		wg.Done()
	}()
	go func() {
		wg.Done()
	}()

	wg.Done() // This should make the counter zero

	select {
	case <-wg.WaitChan():
		t.Fatal("WaitChan should not be closed yet with 1 pending")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}

	wg.Ready()

	select {
	case <-wg.WaitChan():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitChan was not closed after all Dones were called")
	}
}

func TestWaitGroupConcurrentDone(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	count := 100
	wg.Add(count)
	wg.Ready()

	var doneCount int32
	var mu sync.Mutex

	for i := range count {
		go func() {
			time.Sleep(time.Duration(i%5) * time.Millisecond) // Introduce some variation
			mu.Lock()
			doneCount++
			mu.Unlock()
			wg.Done()
		}()
	}

	select {
	case <-wg.WaitChan():
		// Expected
	case <-time.After(500 * time.Millisecond): // Give ample time for all goroutines
		t.Fatal("WaitChan was not closed after concurrent Dones")
	}

	mu.Lock()
	if doneCount != int32(count) {
		t.Errorf("Expected %d Dones, got %d", count, doneCount)
	}
	mu.Unlock()
}

func TestWaitGroupMoreDoneThanAdd(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	wg.Add(1)
	wg.Ready()
	wg.Done()
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected panic when calling Done more times than Add")
		} else if msg := r.(string); msg != alreadyClosedPanic {
			t.Errorf("Unexpected panic message: %s, expected %s", msg, alreadyClosedPanic)
		}
	}()
	wg.Done() // This should panic
}

func TestWaitGroupDoneOnEmptyGroupPanics(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected panic when calling Done on an empty group (no initial Add)")
		} else if msg := r.(string); msg != tooManyDoneCallsPanic {
			t.Errorf("Unexpected panic message: %s, expected %s", msg, tooManyDoneCallsPanic)
		}
	}()
	wg.Done() // This should panic (equivalent to Add(-1) when counter is 0)
}

func TestWaitGroupAddNegativePanic(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected panic when Add(-x) brings counter beneath zero")
		} else if msg := r.(string); msg != tooManyDoneCallsPanic {
			t.Errorf("Unexpected panic message: %s, expected %s", msg, tooManyDoneCallsPanic)
		}
	}()
	wg.Add(1)
	wg.Add(-2)
}

func TestWaitGroupWaitWithoutReady(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup

	select {
	case <-wg.WaitChan():
		t.Fatal("Wait completed despite no work added")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}
}

func TestWaitGroupReadyNoWork(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup

	wg.Ready()

	select {
	case <-wg.WaitChan():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait did not complete despite ready and not work")
	}
}

func TestWaitGroupZeroAddNoCompletion(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	wg.Add(0)

	select {
	case <-wg.WaitChan():
		t.Error("WaitChan should not be closed after Add(0)")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}
}

func TestWaitGroupReuseAfterCompletion(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	wg.Add(1)
	wg.Done()
	wg.Ready()

	select {
	case <-wg.WaitChan():
		// Expected
	case <-time.After(10 * time.Millisecond):
		t.Fatal("WaitChan should be closed")
	}

	// Attempting to reuse should panic
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected panic when trying to reuse a completed WaitGroup")
		} else if msg := r.(string); msg != alreadyClosedPanic {
			t.Errorf("Unexpected panic message: %s, expected %s", msg, alreadyClosedPanic)
		}
	}()
	wg.Add(1)
}

func TestWaitGroupWaitChanMultipleCalls(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	wg.Add(1)
	wg.Ready()

	ch1 := wg.WaitChan()
	ch2 := wg.WaitChan()

	if ch1 != ch2 {
		t.Error("WaitChan should return the same channel instance")
	}

	go func() {
		<-ch1
	}()
	go func() {
		<-ch2
	}()

	wg.Done()

	select {
	case <-ch1:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("ch1 was not closed")
	}
	select {
	case <-ch2:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("ch2 was not closed")
	}
}
func TestWaitGroupWaitChanNestedGoroutines(t *testing.T) {
	t.Parallel()
	var wg chanwg.WaitGroup
	wg.Add(1)
	wg.Ready()
	go func() {
		defer wg.Done()

		wg.Add(1)
		go func() {
			defer wg.Done()
		}()
	}()

	select {
	case <-wg.WaitChan():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("wg was not closed")
	}
}

func TestWaitGroupGo(t *testing.T) {
	t.Parallel()

	var wg chanwg.WaitGroup
	wg.Go(func() {
		wg.Go(func() {
		})
		wg.Go(func() {
			wg.Go(func() {
			})
		})
	})
	wg.Go(func() {
	})
	wg.Ready()

	select {
	case <-wg.WaitChan():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("wg was not closed")
	}
}

func TestWaitGroupGoAllStart(t *testing.T) {
	t.Parallel()
	var counter atomic.Int32

	var wg chanwg.WaitGroup
	wg.Go(func() {
		counter.Add(1)
		wg.Go(func() {
			counter.Add(1)
		})
		wg.Go(func() {
			counter.Add(1)
			wg.Go(func() {
				counter.Add(1)
			})
		})
	})
	wg.Add(1)
	go func() {
		counter.Add(1)
		defer wg.Done()
	}()
	wg.Ready()

	select {
	case <-wg.WaitChan():
		if counter.Load() != 5 {
			t.Errorf("Not all goroutines started, expected 5, got %d", counter.Load())
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("wg was not closed")
	}
}

func TestWaitGroupRace(t *testing.T) {
	t.Parallel()
	timeout := time.After(100 * time.Millisecond)

	for i := 0; i < 1000; i++ {
		var wg chanwg.WaitGroup
		var counter atomic.Int32
		wg.Go(func() {
			counter.Add(1)
		})
		wg.Go(func() {
			counter.Add(1)
		})

		wg.Ready()
		select {
		case <-wg.WaitChan():
		case <-timeout:
			t.Fatal("wg was not closed, is it late?")
		}
		if counter.Load() != 2 {
			t.Fatal("WaitChan closed before all goroutines completed")
		}
	}
}

func BenchmarkChannelWaitGroup(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg chanwg.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
			}()
			wg.Ready()
			<-wg.WaitChan()
		}
	})
}
func BenchmarkSyncWaitGroup(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				wg.Done()
			}()
			wg.Wait()
		}
	})
}
