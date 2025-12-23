# Concurrency in Go

## Table of Contents
1. [Goroutines](#goroutines)
2. [Channels](#channels)
3. [Context](#context)
4. [Sync Primitives](#sync-primitives)
5. [Common Patterns](#common-patterns)
6. [Pitfalls](#pitfalls)

## Goroutines

### Basic Usage

```go
// Launch goroutine
go func() {
    // Do work
}()

// With parameter capture
for _, item := range items {
    item := item // Capture loop variable!
    go func() {
        process(item)
    }()
}

// Or pass as parameter
for _, item := range items {
    go func(i Item) {
        process(i)
    }(item)
}
```

### Waiting for Goroutines

```go
func ProcessAll(items []Item) {
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        item := item
        go func() {
            defer wg.Done()
            process(item)
        }()
    }

    wg.Wait() // Block until all done
}
```

## Channels

### Basic Channel Operations

```go
// Unbuffered channel - synchronous
ch := make(chan int)

// Buffered channel - async up to buffer size
ch := make(chan int, 100)

// Send
ch <- value

// Receive
value := <-ch

// Receive with ok (check if closed)
value, ok := <-ch
if !ok {
    // Channel closed
}

// Close (only sender should close)
close(ch)
```

### Channel Patterns

```go
// Generator pattern
func Generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// Fan-out: multiple goroutines reading from same channel
func FanOut(in <-chan int, workers int) []<-chan int {
    outs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        outs[i] = worker(in)
    }
    return outs
}

// Fan-in: merge multiple channels into one
func FanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    for _, ch := range channels {
        wg.Add(1)
        ch := ch
        go func() {
            defer wg.Done()
            for v := range ch {
                out <- v
            }
        }()
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### Select Statement

```go
func Process(ctx context.Context, in <-chan int) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case v, ok := <-in:
            if !ok {
                return nil // Channel closed
            }
            handle(v)
        }
    }
}

// With timeout
select {
case result := <-resultCh:
    return result, nil
case <-time.After(5 * time.Second):
    return nil, ErrTimeout
}

// Non-blocking check
select {
case v := <-ch:
    // Got value
default:
    // Channel empty, don't block
}
```

## Context

### Context Basics

```go
// Create root context
ctx := context.Background()

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// With deadline
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// With value (use sparingly)
ctx = context.WithValue(ctx, "requestID", "abc123")
```

### Using Context for Cancellation

```go
func LongOperation(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err() // context.Canceled or context.DeadlineExceeded
        default:
            // Do work chunk
            if done := doWorkChunk(); done {
                return nil
            }
        }
    }
}

// Pass context to child operations
func ProcessItems(ctx context.Context, items []Item) error {
    for _, item := range items {
        if err := ctx.Err(); err != nil {
            return err // Early exit if canceled
        }
        if err := processItem(ctx, item); err != nil {
            return err
        }
    }
    return nil
}
```

## Sync Primitives

### Mutex

```go
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

### RWMutex

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]Item
}

func (c *Cache) Get(key string) (Item, bool) {
    c.mu.RLock() // Read lock - multiple readers allowed
    defer c.mu.RUnlock()
    item, ok := c.items[key]
    return item, ok
}

func (c *Cache) Set(key string, item Item) {
    c.mu.Lock() // Write lock - exclusive
    defer c.mu.Unlock()
    c.items[key] = item
}
```

### Once

```go
var (
    instance *Service
    once     sync.Once
)

func GetService() *Service {
    once.Do(func() {
        instance = &Service{}
        instance.Initialize()
    })
    return instance
}
```

### sync.Map (for concurrent map access)

```go
var cache sync.Map

// Store
cache.Store("key", value)

// Load
if v, ok := cache.Load("key"); ok {
    // Use v
}

// LoadOrStore
actual, loaded := cache.LoadOrStore("key", newValue)

// Delete
cache.Delete("key")

// Range
cache.Range(func(key, value any) bool {
    // Process each entry
    return true // Continue iteration
})
```

## Common Patterns

### Worker Pool

```go
func WorkerPool(ctx context.Context, jobs <-chan Job, workers int) <-chan Result {
    results := make(chan Result)

    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                case job, ok := <-jobs:
                    if !ok {
                        return
                    }
                    results <- process(job)
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

### errgroup for Parallel Tasks

```go
import "golang.org/x/sync/errgroup"

func ProcessParallel(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)

    for _, item := range items {
        item := item // Capture
        g.Go(func() error {
            return process(ctx, item)
        })
    }

    return g.Wait() // Returns first error, cancels others
}

// With concurrency limit
func ProcessWithLimit(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(10) // Max 10 concurrent

    for _, item := range items {
        item := item
        g.Go(func() error {
            return process(ctx, item)
        })
    }

    return g.Wait()
}
```

### Semaphore

```go
import "golang.org/x/sync/semaphore"

var sem = semaphore.NewWeighted(10) // Max 10 concurrent

func DoWork(ctx context.Context) error {
    if err := sem.Acquire(ctx, 1); err != nil {
        return err
    }
    defer sem.Release(1)

    // Do limited work
    return nil
}
```

## Pitfalls

### Loop Variable Capture

```go
// BUG - all goroutines share same 'i'
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i) // Prints 10 ten times!
    }()
}

// FIX 1 - capture in new variable
for i := 0; i < 10; i++ {
    i := i
    go func() {
        fmt.Println(i)
    }()
}

// FIX 2 - pass as parameter
for i := 0; i < 10; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}
```

### Goroutine Leaks

```go
// BUG - goroutine blocks forever if no receiver
func leak() {
    ch := make(chan int)
    go func() {
        ch <- 42 // Blocks forever!
    }()
    // Function returns, goroutine leaked
}

// FIX - use context for cancellation
func noLeak(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case ch <- 42:
        case <-ctx.Done():
            return // Exits if context canceled
        }
    }()
}
```

### Data Races

```go
// BUG - data race
var counter int
for i := 0; i < 1000; i++ {
    go func() {
        counter++ // Race!
    }()
}

// FIX - use atomic
var counter int64
for i := 0; i < 1000; i++ {
    go func() {
        atomic.AddInt64(&counter, 1)
    }()
}

// FIX - use mutex
var mu sync.Mutex
var counter int
for i := 0; i < 1000; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}
```

### Detect with Race Detector

```bash
go test -race ./...
go run -race main.go
```
