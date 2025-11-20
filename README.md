[![Go Reference](https://pkg.go.dev/badge/github.com/MatthiasKunnen/chanwg.svg)](https://pkg.go.dev/github.com/MatthiasKunnen/chanwg)

# chanwg
The Go `chanwg` project contains a channel-based, cancelable, alternative to
[`sync.WaitGroup`](https://pkg.go.dev/sync#WaitGroup).

## Why?
When waiting for a `sync.WaitGroup` with `wg.Wait()`, the goroutine that does the waiting will block
indefinitely until the `WaitGroup` completes. There is no way to cancel the waiting.

Using a separate goroutine to wait, leaks a goroutine and complicates the code:

```go
import (
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	done := make(chan struct{})

	go func() {
		// This goroutine leaks due to a forever wait
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
```

Using `chanwg`'s WaitGroup, this can be achieved without as follows:

```go
import (
	"github.com/MatthiasKunnen/chanwg"
	"time"
)

func main()  {
    ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
    defer cancelFunc()

    var wg chanwg.WaitGroup
    wg.Go(func() {
        // Long-running task
    })
    wg.Ready()

    select {
    case <-wg.WaitChan():
    case <-ctx.Done():
        fmt.Println("Abort")
    }
}
```

## Difference from `sync.WaitGroup`

### `chanwg.WaitGroup` is single use
Instead of reusing it after completion, create a new one.

### `chanwg.WaitGroup` will never complete until `Ready` is called
This allows receiving from the wait channel before adding the work.
