[![Go Reference](https://pkg.go.dev/badge/github.com/MatthiasKunnen/chanwg.svg)](https://pkg.go.dev/github.com/MatthiasKunnen/chanwg)

# chanwg
The Go `chanwg` project contains a channel-based, cancelable, alternative to
[`sync.WaitGroup`](https://pkg.go.dev/sync#WaitGroup).

## Why?
When waiting for a `sync.WaitGroup` with `wg.Wait()`, the goroutine that does the waiting will block
indefinitely until the `WaitGroup` completes. There is no way to cancel the waiting.

Using a separate goroutine to wait, leaks a goroutine:

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

Using `chanwg`'s WaitGroup, this can be achieved without leaking any goroutines as follows:

```go
import (
	"github.com/MatthiasKunnen/chanwg"
	"time"
)

func main()  {
	var wg chanwg.WaitGroup
	wg.Add(1)

	select {
	case <-wg.WaitChan():
	case <-time.After(time.Second):
	}
}
```

## Difference from `sync.WaitGroup`

### `chanwg.WaitGroup` is single use
Instead of reusing it after completion, create a new one.

### `chanwg.WaitGroup` will never complete until at least one `Add` and `Done` is performed
This is done to allow this:
```go
type Foo struct {
	startedWg chanwg.WaitGroup
	started <-chan struct{}
}

func NewFoo() *Foo {
	foo := &Foo{
	}
	foo.started = foo.startedWg.WaitChan()

	return foo
}

func (f *Foo) Start(ctx context.Context) error {
	f.startedWg.Add(1)

	select {
	case <-f.started:
		return nil
	case <-ctx.Done():
        return ctx.Err()
	}
}
```
