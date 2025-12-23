# Performance in Go

## Table of Contents
1. [Profiling](#profiling)
2. [Memory Optimization](#memory-optimization)
3. [CPU Optimization](#cpu-optimization)
4. [Benchmarking](#benchmarking)
5. [Common Optimizations](#common-optimizations)

## Profiling

### CPU Profiling

```go
import "runtime/pprof"

func main() {
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Your code here
}
```

```bash
# Analyze profile
go tool pprof cpu.prof

# Interactive commands
(pprof) top10          # Top 10 functions
(pprof) list funcName  # Source annotation
(pprof) web            # Visualize in browser
```

### Memory Profiling

```go
import "runtime/pprof"

func main() {
    // ... your code ...

    f, _ := os.Create("mem.prof")
    pprof.WriteHeapProfile(f)
    f.Close()
}
```

```bash
go tool pprof mem.prof
(pprof) top10 -inuse_space  # Current memory
(pprof) top10 -alloc_space  # Total allocations
```

### HTTP Profiling Server

```go
import _ "net/http/pprof"

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()

    // Your application
}
```

```bash
# CPU profile for 30 seconds
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### Trace

```go
import "runtime/trace"

func main() {
    f, _ := os.Create("trace.out")
    trace.Start(f)
    defer trace.Stop()

    // Your code
}
```

```bash
go tool trace trace.out
```

## Memory Optimization

### Pre-allocate Slices

```go
// BAD - grows slice repeatedly
var result []int
for i := 0; i < 10000; i++ {
    result = append(result, i) // Multiple allocations
}

// GOOD - pre-allocate
result := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    result = append(result, i) // No reallocation
}

// EVEN BETTER if size known
result := make([]int, 10000)
for i := 0; i < 10000; i++ {
    result[i] = i
}
```

### Avoid String Concatenation in Loops

```go
// BAD - O(n²) allocations
var s string
for _, item := range items {
    s += item.Name + ", "
}

// GOOD - O(n) with builder
var b strings.Builder
for _, item := range items {
    b.WriteString(item.Name)
    b.WriteString(", ")
}
s := b.String()
```

### Sync.Pool for Temporary Objects

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func Process(data []byte) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()

    // Use buf...
    return buf.String()
}
```

### Avoid Escape to Heap

```go
// Check escape analysis
// go build -gcflags="-m" .

// BAD - escapes to heap
func NewUser(name string) *User {
    u := User{Name: name}
    return &u // Escapes
}

// OK if caller needs pointer
// But avoid if not necessary

// Sometimes stack is better
func ProcessUser(name string) Result {
    u := User{Name: name} // Stays on stack
    return process(u)
}
```

### Reduce Allocations with Value Receivers

```go
// Pointer receiver - may cause escape
func (u *User) Name() string {
    return u.name
}

// Value receiver - may stay on stack
func (u User) Name() string {
    return u.name
}

// Use value receivers for small, read-only methods
// Use pointer receivers for large structs or mutations
```

## CPU Optimization

### Avoid Interface Overhead in Hot Paths

```go
// SLOW - interface dispatch
func Sum(values []interface{}) int {
    var sum int
    for _, v := range values {
        sum += v.(int)
    }
    return sum
}

// FAST - direct type
func Sum(values []int) int {
    var sum int
    for _, v := range values {
        sum += v
    }
    return sum
}
```

### Use Efficient Data Structures

```go
// Map lookup O(1) vs slice search O(n)
// Use map for lookups
seen := make(map[string]bool)
for _, item := range items {
    if seen[item.ID] {
        continue
    }
    seen[item.ID] = true
    // process
}
```

### Batch Operations

```go
// SLOW - individual inserts
for _, user := range users {
    db.Insert(user)
}

// FAST - batch insert
db.InsertBatch(users)
```

## Benchmarking

### Write Effective Benchmarks

```go
func BenchmarkProcess(b *testing.B) {
    // Setup (not counted)
    data := setupTestData()

    b.ResetTimer() // Start timing

    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

// With allocations
func BenchmarkProcess(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

// Compare implementations
func BenchmarkOld(b *testing.B) { /* ... */ }
func BenchmarkNew(b *testing.B) { /* ... */ }
```

```bash
# Run benchmarks
go test -bench=. -benchmem ./...

# Compare with benchstat
go test -bench=. -count=10 > old.txt
# Make changes
go test -bench=. -count=10 > new.txt
benchstat old.txt new.txt
```

### Avoid Compiler Optimizations in Benchmarks

```go
var result int // Package-level to prevent optimization

func BenchmarkCompute(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = Compute(i)
    }
    result = r // Prevent dead code elimination
}
```

## Common Optimizations

### JSON Performance

```go
// Standard library - convenient but slower
json.Marshal(data)
json.Unmarshal(data, &v)

// For hot paths, consider:
// - jsoniter: drop-in replacement, faster
// - easyjson: code generation, fastest
// - avoid reflection with custom marshalers
```

### String to Bytes Without Copy

```go
import "unsafe"

// Zero-copy string to bytes (read-only!)
func stringToBytes(s string) []byte {
    return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Usually just use []byte(s) unless profiling shows issue
```

### Reduce GC Pressure

```go
// Reuse objects
type Processor struct {
    buffer []byte // Reused between calls
}

func (p *Processor) Process(data []byte) Result {
    p.buffer = p.buffer[:0] // Reset, keep capacity
    // Use p.buffer...
}

// Use value types for small data
type Point struct { X, Y float64 } // 16 bytes, value type ok

// Use pointers for large data
type LargeStruct struct { /* many fields */ }
func Process(s *LargeStruct) { /* ... */ }
```

### Parallel Processing

```go
import "runtime"

func ProcessParallel(items []Item) []Result {
    numWorkers := runtime.GOMAXPROCS(0)
    results := make([]Result, len(items))

    var wg sync.WaitGroup
    chunkSize := (len(items) + numWorkers - 1) / numWorkers

    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := min(start+chunkSize, len(items))
        if start >= end {
            break
        }

        wg.Add(1)
        go func(start, end int) {
            defer wg.Done()
            for j := start; j < end; j++ {
                results[j] = process(items[j])
            }
        }(start, end)
    }

    wg.Wait()
    return results
}
```
