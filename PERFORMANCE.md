# Performance Tuning Tips for Go

A comprehensive guide to performance optimization techniques in Go applications.

## Table of Contents

- [Profiling and Benchmarking](#profiling-and-benchmarking)
- [Memory Management](#memory-management)
- [Concurrency](#concurrency)
- [Data Structures](#data-structures)
- [I/O Operations](#io-operations)
- [Compilation](#compilation)
- [Common Anti-Patterns](#common-anti-patterns)

---

## Profiling and Benchmarking

### 1. **Use Built-in Profiling Tools**

Always profile before optimizing. Go provides excellent profiling tools:

```go
import _ "net/http/pprof"

// In your main function
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

Access profiles at `http://localhost:6060/debug/pprof/`

### 2. **Run Benchmarks with Memory Stats**

```bash
go test -bench=. -benchmem
go test -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out
```

### 3. **Use `go tool pprof` for Analysis**

```bash
go tool pprof -http=:8080 cpu.out
go tool pprof -http=:8080 -alloc_space mem.out
```

### 4. **Trace Execution**

```bash
go test -trace=trace.out
go tool trace trace.out
```

### 5. **Benchmark Comparison**

Use `benchstat` to compare benchmarks statistically:

```bash
go get -tool golang.org/x/perf/cmd/benchstat@latest
go tool benchstat old.txt new.txt
```

---

## Memory Management

### 6. **Pre-allocate Slices**

```go
// Bad
var items []Item
for i := 0; i < n; i++ {
    items = append(items, Item{})  // Multiple allocations
}

// Good
items := make([]Item, 0, n)  // Single allocation
for i := 0; i < n; i++ {
    items = append(items, Item{})
}

// Best (if you know all values upfront)
items := make([]Item, n)
for i := 0; i < n; i++ {
    items[i] = Item{}
}
```

### 7. **Reuse Buffers with sync.Pool**

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process() {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    buf.Reset()
    
    // Use buffer
}
```

### 8. **Avoid String Concatenation in Loops**

```go
// Bad
var s string
for i := 0; i < n; i++ {
    s += "item"  // Creates new string each iteration
}

// Good
var sb strings.Builder
sb.Grow(n * 4)  // Preallocate if size is known
for i := 0; i < n; i++ {
    sb.WriteString("item")
}
s := sb.String()
```

### 9. **Understand Escape Analysis**

```bash
go build -gcflags='-m -l'  # See what escapes to heap
```

Keep variables on the stack when possible:

```go
// Bad - escapes to heap
func createUser() *User {
    u := User{Name: "John"}
    return &u  // Escapes
}

// Good - stays on stack (if User is small)
func createUser() User {
    return User{Name: "John"}
}
```

### 10. **Use Pointer Receivers Appropriately**

```go
// Use pointers for large structs or when you need to modify
func (u *User) Update(name string) {  // Good
    u.Name = name
}

// Use values for small structs (< 64 bytes typically)
func (p Point) Distance() float64 {  // Good
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}
```

### 11. **Reduce Allocations in Hot Paths**

```go
// Bad
func process(data []byte) {
    s := string(data)  // Allocation
    // ...
}

// Good - avoid conversion if possible
func process(data []byte) {
    // Work with []byte directly
}
```

---

## Concurrency

### 12. **Use Worker Pools**

```go
func processItems(items []Item, numWorkers int) {
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))
    
    // Start workers
    for w := 0; w < numWorkers; w++ {
        go worker(jobs, results)
    }
    
    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // Collect results
    for i := 0; i < len(items); i++ {
        <-results
    }
}
```

### 13. **Avoid Goroutine Leaks**

```go
// Bad - goroutine may leak
func process() {
    ch := make(chan int)
    go func() {
        ch <- compute()  // Blocks forever if nobody reads
    }()
}

// Good - use context or buffered channel
func process(ctx context.Context) {
    ch := make(chan int, 1)  // Buffered
    go func() {
        select {
        case ch <- compute():
        case <-ctx.Done():
            return
        }
    }()
}
```

### 14. **Right-Size Buffered Channels**

```go
// Unbuffered - causes blocking
ch := make(chan int)

// Buffered - reduces blocking
ch := make(chan int, 100)

// Match producer/consumer rate
ch := make(chan int, runtime.NumCPU())
```

### 15. **Use sync.Once for Initialization**

```go
var (
    instance *DB
    once     sync.Once
)

func GetDB() *DB {
    once.Do(func() {
        instance = &DB{}
        instance.Connect()
    })
    return instance
}
```

### 16. **Avoid Mutex Contention**

```go
// Bad - single mutex for all operations
type Cache struct {
    sync.Mutex
    items map[string]interface{}
}

// Good - use sync.Map or RWMutex
type Cache struct {
    sync.RWMutex  // Allows multiple readers
    items map[string]interface{}
}

// Or use sync.Map for concurrent access
var cache sync.Map
```

### 17. **Consider GOMAXPROCS**

```go
import "runtime"

func init() {
    // Usually default is fine (number of CPUs)
    // But you can adjust for CPU-bound vs I/O-bound
    runtime.GOMAXPROCS(runtime.NumCPU())
}
```

### 18. **Use sync.WaitGroup Properly**

Coordinate multiple goroutines without leaks:

```go
// Bad - race condition
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    go func() {
        wg.Add(1)  // Wrong! Can cause race
        defer wg.Done()
        // work...
    }()
}
wg.Wait()

// Good - Add before goroutine starts
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)  // Add in parent goroutine
    go func() {
        defer wg.Done()
        // work...
    }()
}
wg.Wait()

// Best - batch Add for known count
var wg sync.WaitGroup
wg.Add(10)  // Add once for all goroutines
for i := 0; i < 10; i++ {
    go func() {
        defer wg.Done()
        // work...
    }()
}
wg.Wait()
```

**Key points:**
- Always call `Add()` before starting the goroutine
- Use `defer Done()` to ensure it's always called
- Can batch `Add(n)` for known counts

### 19. **Use sync.Cond for Efficient Waiting**

For scenarios where goroutines need to wait for conditions:

```go
// Bad - busy waiting (wastes CPU)
var ready bool
var mu sync.Mutex

// Goroutine waiting
for {
    mu.Lock()
    if ready {
        mu.Unlock()
        break
    }
    mu.Unlock()
    time.Sleep(10 * time.Millisecond)  // Wasteful
}

// Good - use sync.Cond
var ready bool
var mu sync.Mutex
cond := sync.NewCond(&mu)

// Waiting goroutine
mu.Lock()
for !ready {
    cond.Wait()  // Efficiently blocks until signaled
}
// work with ready condition
mu.Unlock()

// Signaling goroutine
mu.Lock()
ready = true
cond.Signal()  // Wake one waiter
// or cond.Broadcast()  // Wake all waiters
mu.Unlock()
```

**Use cases:**
- Producer-consumer with multiple consumers
- State changes that multiple goroutines wait on
- Implementing custom synchronization primitives

**Example: Bounded queue**

```go
type Queue struct {
    mu    sync.Mutex
    notEmpty *sync.Cond
    notFull  *sync.Cond
    items []interface{}
    capacity int
}

func NewQueue(capacity int) *Queue {
    q := &Queue{
        items: make([]interface{}, 0, capacity),
        capacity: capacity,
    }
    q.notEmpty = sync.NewCond(&q.mu)
    q.notFull = sync.NewCond(&q.mu)
    return q
}

func (q *Queue) Enqueue(item interface{}) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    for len(q.items) == q.capacity {
        q.notFull.Wait()  // Wait until not full
    }
    
    q.items = append(q.items, item)
    q.notEmpty.Signal()  // Signal waiting consumers
}

func (q *Queue) Dequeue() interface{} {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    for len(q.items) == 0 {
        q.notEmpty.Wait()  // Wait until not empty
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    q.notFull.Signal()  // Signal waiting producers
    return item
}
```

---

## Data Structures

### 20. **Choose the Right Map Size**

```go
// Preallocate maps if size is known
m := make(map[string]int, expectedSize)
```

### 21. **Use Struct Keys for Maps Carefully**

```go
// Comparable structs can be used as keys
type Point struct {
    X, Y int
}
m := make(map[Point]string)

// But consider using a simple key instead
m := make(map[string]string)  // key: fmt.Sprintf("%d,%d", x, y)
```

### 22. **Prefer Arrays over Slices for Fixed Size**

```go
// If size is known and small
var buffer [4096]byte  // Stack allocated

// vs
buffer := make([]byte, 4096)  // Heap allocated
```

### 23. **Use Appropriate Integer Types**

```go
// Use smallest type that fits
var count uint8     // 0-255
var id    uint32    // 0-4,294,967,295

// But don't micro-optimize - int is usually fine
var normalSize int
```

---

## I/O Operations

### 24. **Use Buffered I/O**

```go
// Bad
file, _ := os.Open("large.txt")
defer file.Close()
scanner := bufio.NewScanner(file)

// Good - specify buffer size for large files
scanner := bufio.NewScanner(file)
scanner.Buffer(make([]byte, 64*1024), 1024*1024)
```

### 25. **Batch Database Operations**

```go
// Bad - individual inserts
for _, item := range items {
    db.Exec("INSERT INTO ...", item)
}

// Good - batch insert
tx, _ := db.Begin()
stmt, _ := tx.Prepare("INSERT INTO ...")
for _, item := range items {
    stmt.Exec(item)
}
tx.Commit()
```

### 26. **Use Connection Pools**

```go
db, err := sql.Open("postgres", dsn)
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 27. **Avoid Excessive JSON Marshal/Unmarshal**

```go
// Consider using easyjson or jsoniter for performance
// go get github.com/mailru/easyjson

// Or use sync.Pool for encoder/decoder
var encoderPool = sync.Pool{
    New: func() interface{} {
        return json.NewEncoder(nil)
    },
}
```

---

## Compilation

### 28. **Use Build Tags for Conditional Compilation**

```go
// +build !race

// Code that should not be compiled with race detector
```

### 29. **Disable Bounds Checking in Critical Loops**

```go
// After verifying safety
for i := 0; i < len(data); i++ {
    _ = data[i] // Compiler may eliminate bounds check
}
```

### 30. **Use Inlining Hints**

```bash
go build -gcflags='-m'  # See inlining decisions
```

```go
// Small functions (< 80 nodes) are usually inlined
// Use //go:noinline to prevent inlining
//go:noinline
func expensiveSetup() {}
```

### 31. **Enable Compiler Optimizations**

```bash
go build -ldflags="-s -w"  # Strip debug info (smaller binary)
go build -tags netgo       # Pure Go networking (no cgo)
```

---

## Common Anti-Patterns

### 32. **Don't Defer in Loops**

```go
// Bad - defer accumulates
for _, file := range files {
    f, _ := os.Open(file)
    defer f.Close()  // All close at end of function
}

// Good - close explicitly or use a function
for _, file := range files {
    func() {
        f, _ := os.Open(file)
        defer f.Close()
    }()
}
```

### 33. **Avoid Allocating in Hot Paths**

```go
// Bad
func process(data []byte) Result {
    temp := make([]byte, len(data))  // Allocation in hot path
    copy(temp, data)
    // ...
}

// Good - reuse buffer
type Processor struct {
    buffer []byte
}

func (p *Processor) process(data []byte) Result {
    if cap(p.buffer) < len(data) {
        p.buffer = make([]byte, len(data))
    }
    p.buffer = p.buffer[:len(data)]
    copy(p.buffer, data)
    // ...
}
```

### 34. **Don't Ignore Context**

```go
// Bad
func fetchData(url string) ([]byte, error) {
    resp, err := http.Get(url)
    // ...
}

// Good
func fetchData(ctx context.Context, url string) ([]byte, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    // ...
}
```

### 35. **Avoid Time.After in Loops**

```go
// Bad - leaks tickers
for {
    select {
    case <-time.After(1 * time.Second):  // New timer each iteration
        // ...
    }
}

// Good
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for {
    select {
    case <-ticker.C:
        // ...
    }
}
```

### 36. **Use Constants for Repeated Values**

```go
// Bad
if size > 1024*1024 {  // Magic number
    // ...
}

// Good
const OneMB = 1024 * 1024

if size > OneMB {
    // ...
}
```

### 37. **Profile in Production Environments**

```go
// Add conditional profiling
import _ "net/http/pprof"

if os.Getenv("ENABLE_PPROF") == "true" {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

---

## Performance Checklist

Before deploying to production:

- [ ] Run benchmarks and profiling
- [ ] Check for goroutine leaks (`runtime.NumGoroutine()`)
- [ ] Verify memory usage patterns
- [ ] Test under load
- [ ] Monitor garbage collection pauses
- [ ] Review hot paths with profiler
- [ ] Ensure proper connection pooling
- [ ] Check for excessive allocations
- [ ] Verify context usage and timeouts
- [ ] Test with race detector (`go test -race`)

---

## Resources

- [Go Performance Workshop](https://github.com/davecheney/high-performance-go-workshop)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [Go Memory Management](https://go.dev/doc/gc-guide)
- [Diagnostics](https://go.dev/doc/diagnostics)
- [Go Performance Book](https://github.com/dgryski/go-perfbook)

---

**Remember**: Premature optimization is the root of all evil. Always:
1. Write correct code first
2. Profile to find bottlenecks
3. Optimize based on data
4. Measure improvements
