# Complete Guide to pprof Profiling in Go

A comprehensive guide to using Go's built-in profiling tool `pprof` for performance analysis and optimization.

## Table of Contents

- [Introduction](#introduction)
- [Profile Types](#profile-types)
- [Web Application Profiling](#web-application-profiling)
- [CLI Application Profiling](#cli-application-profiling)
- [Analyzing Profiles](#analyzing-profiles)
- [Visualization Techniques](#visualization-techniques)
- [Advanced Usage](#advanced-usage)
- [Best Practices](#best-practices)
- [Common Scenarios](#common-scenarios)

---

## Introduction

`pprof` is Go's built-in profiling tool that helps identify performance bottlenecks in your applications. It can profile:
- CPU usage
- Memory allocations
- Goroutine blocking
- Mutex contention
- Execution traces

### Why Profile?

> **Don't guess, measure!**

Profiling helps you:
- Identify actual bottlenecks (not assumed ones)
- Optimize based on data
- Detect memory leaks
- Find goroutine leaks
- Understand concurrency issues

---

## Profile Types

### 1. **CPU Profile**

Tracks where your program spends CPU time.

**What it shows:**
- Functions consuming the most CPU
- Hot paths in your code
- Recursive call overhead

**When to use:**
- Application is slow
- High CPU usage
- Need to optimize compute-intensive code

**CLI Usage:**
```bash
# Capture CPU profile
go run main.go -cpuprofile=cpu.prof

# Or with go test
go test -cpuprofile=cpu.prof
```

**Web Usage:**
```bash
# Access via HTTP (captures for 30 seconds by default)
curl http://localhost:8080/debug/pprof/profile > cpu.prof

# Shorter duration
curl http://localhost:8080/debug/pprof/profile?seconds=10 > cpu.prof
```

### 2. **Heap Profile**

Shows memory allocations and current heap usage.

**What it shows:**
- Memory allocated by each function
- Current heap objects (in-use)
- Total allocations over time

**When to use:**
- High memory usage
- Suspected memory leaks
- Want to reduce allocations

**CLI Usage:**
```bash
go run main.go -memprofile=mem.prof
```

**Web Usage:**
```bash
# Current heap (in-use objects)
curl http://localhost:8080/debug/pprof/heap > heap.prof

# All allocations (since program start)
curl http://localhost:8080/debug/pprof/allocs > allocs.prof
```

**Key Metrics:**
- `alloc_space`: Total bytes allocated
- `alloc_objects`: Total objects allocated
- `inuse_space`: Currently allocated bytes
- `inuse_objects`: Currently allocated objects

### 3. **Goroutine Profile**

Shows all currently running goroutines and their call stacks.

**What it shows:**
- Number of goroutines
- What each goroutine is doing
- Where goroutines are blocked

**When to use:**
- Goroutine leaks
- Too many goroutines
- Debugging concurrency issues

**CLI Usage:**
```bash
# Programmatically dump goroutines
pprof.Lookup("goroutine").WriteTo(file, 0)
```

**Web Usage:**
```bash
curl http://localhost:8080/debug/pprof/goroutine > goroutine.prof
```

### 4. **Block Profile**

Tracks where goroutines block waiting on synchronization primitives.

**What it shows:**
- Channel operations blocking time
- Mutex lock contention
- Select statement blocking

**When to use:**
- Application seems stuck
- High latency
- Concurrency bottlenecks

**Setup Required:**
```go
import "runtime"

func init() {
    // Enable block profiling (rate = 1 means track all)
    runtime.SetBlockProfileRate(1)
}
```

**CLI Usage:**
```bash
go run main.go -blockprofile=block.prof
```

**Web Usage:**
```bash
curl http://localhost:8080/debug/pprof/block > block.prof
```

### 5. **Mutex Profile**

Tracks contention on mutexes.

**What it shows:**
- Which mutexes have the most contention
- How long goroutines wait for mutexes
- Lock/unlock patterns

**When to use:**
- High mutex contention suspected
- Optimizing concurrent data access
- Scaling issues with goroutines

**Setup Required:**
```go
import "runtime"

func init() {
    // Enable mutex profiling
    // Fraction: 1 = sample every event, 5 = sample 1/5 events
    runtime.SetMutexProfileFraction(1)
}
```

**CLI Usage:**
```bash
go run main.go -mutexprofile=mutex.prof
```

**Web Usage:**
```bash
curl http://localhost:8080/debug/pprof/mutex > mutex.prof
```

### 6. **Execution Trace**

Records detailed execution events over time.

**What it shows:**
- Goroutine creation/blocking/unblocking
- System calls
- GC events
- Processor utilization

**When to use:**
- Understanding goroutine scheduling
- Analyzing concurrency patterns
- Finding subtle race conditions
- GC impact analysis

**CLI Usage:**
```bash
# Import trace package
import "runtime/trace"

# Enable in code
f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()

# Or with flag
go run main.go -trace=trace.out
```

**View trace:**
```bash
go tool trace trace.out
```

### 7. **Allocs Profile**

Tracks all memory allocations (historical).

**What it shows:**
- Every allocation made since program start
- Allocation patterns over time

**When to use:**
- Understanding total allocation behavior
- Finding allocation hot spots
- Reducing GC pressure

**Web Usage:**
```bash
curl http://localhost:8080/debug/pprof/allocs > allocs.prof
```

### 8. **Thread Creation Profile**

Shows OS thread creation.

**What it shows:**
- Stack traces of thread creation

**When to use:**
- Debugging thread-related issues
- Rarely needed in typical applications

**Web Usage:**
```bash
curl http://localhost:8080/debug/pprof/threadcreate > thread.prof
```

---

## Web Application Profiling

### Setup

```go
package main

import (
    "net/http"
    _ "net/http/pprof"  // Import for side effects
    "runtime"
)

func main() {
    // Enable block and mutex profiling
    runtime.SetBlockProfileRate(1)
    runtime.SetMutexProfileFraction(1)
    
    // Your application routes
    http.HandleFunc("/", handler)
    
    // pprof automatically adds handlers at /debug/pprof/
    http.ListenAndServe(":8080", nil)
}
```

### Running the Example

```bash
cd examples/webpprof
go run main.go
```

### Access Profiles

Open browser to http://localhost:8080/debug/pprof/ to see:
- `/debug/pprof/` - Index page
- `/debug/pprof/heap` - Heap profile
- `/debug/pprof/goroutine` - Goroutine profile
- `/debug/pprof/profile` - CPU profile (30s)
- `/debug/pprof/block` - Block profile
- `/debug/pprof/mutex` - Mutex profile
- `/debug/pprof/allocs` - Allocation profile
- `/debug/pprof/threadcreate` - Thread creation

### Generate Load

```bash
# Create users (memory allocations)
curl "http://localhost:8080/api/users?count=1000"

# CPU intensive task
curl "http://localhost:8080/api/compute?iterations=1000000"

# Memory allocations
curl "http://localhost:8080/api/allocate?size=1000"

# Goroutine leak simulation
curl "http://localhost:8080/api/leak?count=50"
```

### Capture Profiles

```bash
# CPU profile (10 seconds)
curl http://localhost:8080/debug/pprof/profile?seconds=10 > cpu.prof

# Heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof

# Goroutine profile
curl http://localhost:8080/debug/pprof/goroutine > goroutine.prof

# Block profile
curl http://localhost:8080/debug/pprof/block > block.prof

# Mutex profile
curl http://localhost:8080/debug/pprof/mutex > mutex.prof
```

---

## CLI Application Profiling

### Running the Example

```bash
cd examples/clipprof

# CPU profiling with all workload
go run main.go -cpuprofile=cpu.prof -workload=all -duration=10

# Memory profiling
go run main.go -memprofile=mem.prof -workload=memory -allocsize=500

# Multiple profiles at once
go run main.go \
  -cpuprofile=cpu.prof \
  -memprofile=mem.prof \
  -blockprofile=block.prof \
  -mutexprofile=mutex.prof \
  -trace=trace.out \
  -workload=all \
  -duration=10

# Specific workload types
go run main.go -cpuprofile=cpu.prof -workload=cpu -iterations=10000
go run main.go -memprofile=mem.prof -workload=memory -allocsize=1000
go run main.go -trace=trace.out -workload=goroutines -goroutines=200
```

### CLI Flags

- `-cpuprofile=file` - Enable CPU profiling
- `-memprofile=file` - Enable memory profiling
- `-blockprofile=file` - Enable block profiling
- `-mutexprofile=file` - Enable mutex profiling
- `-trace=file` - Enable execution trace
- `-workload=type` - cpu, memory, goroutines, or all
- `-iterations=N` - Iterations for CPU workload
- `-allocsize=MB` - Size in MB for memory workload
- `-goroutines=N` - Number of goroutines to spawn
- `-duration=seconds` - Duration to run workload

---

## Analyzing Profiles

### Interactive Mode

Most powerful way to analyze profiles:

```bash
go tool pprof cpu.prof
```

**Commands in interactive mode:**
```
(pprof) top          # Show top functions by resource usage
(pprof) top10        # Show top 10 functions
(pprof) top -cum     # Sort by cumulative time
(pprof) list main.   # Show source code for main package functions
(pprof) list main.computeFibonacci  # Show specific function
(pprof) web          # Open graphviz visualization
(pprof) weblist main.processData   # Show annotated source in browser
(pprof) peek fibonacci              # Show callers and callees
(pprof) traces       # Show sample traces
(pprof) help         # Show all commands
```

### Web UI (Recommended)

```bash
# Start web server on port 8080
go tool pprof -http=:8080 cpu.prof

# Compare profiles
go tool pprof -http=:8080 -base=before.prof after.prof

# Multiple profiles
go tool pprof -http=:8080 cpu.prof mem.prof
```

**Web UI Features:**
- Interactive flame graphs
- Call graphs
- Source code view
- Top functions table
- Comparison mode

### Command Line Analysis

```bash
# Top 20 functions
go tool pprof -top cpu.prof

# Show cumulative times
go tool pprof -top -cum cpu.prof

# List specific function
go tool pprof -list=main.processData cpu.prof

# Filter by focus
go tool pprof -top -focus=processData cpu.prof

# Ignore certain functions
go tool pprof -top -ignore=runtime cpu.prof

# Text output
go tool pprof -text cpu.prof > analysis.txt
```

### Memory Profile Analysis

```bash
# In-use memory (current heap)
go tool pprof -http=:8080 -sample_index=inuse_space heap.prof

# Total allocations (historical)
go tool pprof -http=:8080 -sample_index=alloc_space heap.prof

# Count of objects
go tool pprof -http=:8080 -sample_index=inuse_objects heap.prof
go tool pprof -http=:8080 -sample_index=alloc_objects heap.prof
```

**Sample indices for heap:**
- `inuse_space` - Memory currently in use (bytes)
- `inuse_objects` - Objects currently in use (count)
- `alloc_space` - Total allocated memory (bytes)
- `alloc_objects` - Total allocated objects (count)

### Goroutine Profile Analysis

```bash
# Interactive analysis
go tool pprof goroutine.prof
(pprof) top
(pprof) traces

# Web view
go tool pprof -http=:8080 goroutine.prof

# Text dump
go tool pprof -text goroutine.prof
```

### Block Profile Analysis

```bash
# Web view (recommended)
go tool pprof -http=:8080 block.prof

# Show where blocking occurs
go tool pprof -top block.prof

# View in nanoseconds (default)
go tool pprof -unit=ns -top block.prof

# View in milliseconds
go tool pprof -unit=ms -top block.prof
```

---

## Visualization Techniques

### 1. Flame Graph

Best for understanding call hierarchies and hot paths.

```bash
go tool pprof -http=:8080 cpu.prof
# Click "Flame Graph" in web UI
```

**How to read:**
- Width = resource usage (wider = more)
- Height = call stack depth
- Color = different functions
- Click to zoom in

### 2. Graph View

Shows call relationships with metrics.

```bash
go tool pprof -http=:8080 cpu.prof
# Default view is graph
```

**How to read:**
- Box size = resource usage
- Arrows = call relationships
- Numbers = actual metrics
- Red = hot paths

### 3. Source View

Shows annotated source code with metrics.

```bash
go tool pprof -http=:8080 cpu.prof
# Click "Source" tab
```

### 4. Top Table

List of top resource consumers.

```bash
go tool pprof -http=:8080 cpu.prof
# Click "Top" tab
```

### 5. Peek View

Shows callers and callees of a function.

```bash
go tool pprof cpu.prof
(pprof) peek processData
```

### 6. Command Line Visualization

```bash
# Generate SVG
go tool pprof -svg cpu.prof > cpu.svg

# Generate PNG
go tool pprof -png cpu.prof > cpu.png

# Generate PDF
go tool pprof -pdf cpu.prof > cpu.pdf

# Generate DOT format
go tool pprof -dot cpu.prof > cpu.dot
```

---

## Advanced Usage

### Comparing Profiles

Compare before and after optimization:

```bash
# Capture baseline
curl http://localhost:8080/debug/pprof/profile?seconds=30 > before.prof

# Make changes, restart, generate load

# Capture new profile
curl http://localhost:8080/debug/pprof/profile?seconds=30 > after.prof

# Compare
go tool pprof -http=:8080 -base=before.prof after.prof

# Command line diff
go tool pprof -base=before.prof after.prof
(pprof) top -cum
```

### Filtering and Focusing

```bash
# Focus on specific functions
go tool pprof -focus=processData cpu.prof

# Ignore runtime functions
go tool pprof -ignore=runtime cpu.prof

# Show only paths through a function
go tool pprof -focus=main.* -ignore=runtime.* cpu.prof

# Use regex
go tool pprof -focus="(process|compute)" cpu.prof
```

### Custom Sample Index

```bash
# For memory profiles
go tool pprof -sample_index=alloc_space heap.prof
go tool pprof -sample_index=inuse_space heap.prof
go tool pprof -sample_index=alloc_objects heap.prof
go tool pprof -sample_index=inuse_objects heap.prof
```

### Different Units

```bash
# Show in different units
go tool pprof -unit=ms block.prof   # milliseconds
go tool pprof -unit=us block.prof   # microseconds
go tool pprof -unit=ns block.prof   # nanoseconds
go tool pprof -unit=s cpu.prof      # seconds
go tool pprof -unit=MB heap.prof    # megabytes
```

### Remote Profiling

```bash
# Profile remote server directly
go tool pprof http://production.example.com/debug/pprof/profile

# With duration
go tool pprof http://production.example.com/debug/pprof/profile?seconds=60

# Heap from remote
go tool pprof http://production.example.com/debug/pprof/heap
```

### Continuous Profiling

```bash
# Capture multiple profiles over time
for i in {1..10}; do
    curl http://localhost:8080/debug/pprof/profile?seconds=10 > cpu_$i.prof
    sleep 60
done

# Analyze trends
go tool pprof -http=:8080 cpu_*.prof
```

### Programmatic Profile Writing

```go
import (
    "os"
    "runtime/pprof"
)

// Write heap profile
func writeHeapProfile(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()
    
    return pprof.WriteHeapProfile(f)
}

// Write any profile type
func writeProfile(profileName, filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()
    
    p := pprof.Lookup(profileName)
    return p.WriteTo(f, 0)
}

// Available profiles:
// - "goroutine"
// - "heap"
// - "allocs"
// - "threadcreate"
// - "block"
// - "mutex"
```

---

## Best Practices

### 1. Profile Production Systems Safely

```go
// Only enable pprof in production with authentication
if os.Getenv("ENABLE_PPROF") == "true" {
    go func() {
        // Bind to localhost only
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

### 2. Use Appropriate Sample Rates

```go
// Block profiling
runtime.SetBlockProfileRate(1) // Development: capture all
runtime.SetBlockProfileRate(10000) // Production: sample 1 per 10,000 ns

// Mutex profiling
runtime.SetMutexProfileFraction(1)    // Development: capture all
runtime.SetMutexProfileFraction(1000) // Production: sample 1/1000 events
```

### 3. Profile Long Enough

```bash
# Too short - might miss patterns
curl http://localhost:8080/debug/pprof/profile?seconds=5

# Good - captures meaningful data
curl http://localhost:8080/debug/pprof/profile?seconds=30

# For variable workloads
curl http://localhost:8080/debug/pprof/profile?seconds=60
```

### 4. Generate Realistic Load

```bash
# Use load testing tools
hey -n 10000 -c 100 http://localhost:8080/api/process

# Or ab (Apache Bench)
ab -n 10000 -c 100 http://localhost:8080/api/process

# While profiling
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof &
hey -n 10000 -c 100 http://localhost:8080/api/process
```

### 5. Compare Apples to Apples

```bash
# Same workload, same duration, same conditions
curl http://localhost:8080/debug/pprof/profile?seconds=30 > before.prof

# Make changes

curl http://localhost:8080/debug/pprof/profile?seconds=30 > after.prof
go tool pprof -base=before.prof after.prof
```

### 6. Focus on Hotspots

Look for functions that are:
- High in `flat` (time spent in function itself)
- High in `cum` (cumulative time including callees)
- Called frequently
- In your code (not runtime/stdlib)

### 7. Combine Multiple Profile Types

```bash
# Get complete picture
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof
curl http://localhost:8080/debug/pprof/heap > heap.prof
curl http://localhost:8080/debug/pprof/goroutine > goroutine.prof
curl http://localhost:8080/debug/pprof/block > block.prof

# Analyze together
go tool pprof -http=:8080 cpu.prof heap.prof
```

---

## Common Scenarios

### Scenario 1: High CPU Usage

```bash
# 1. Capture CPU profile under load
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof

# 2. Analyze
go tool pprof -http=:8080 cpu.prof

# 3. Look for:
#    - Recursive functions (deep call stacks)
#    - Tight loops
#    - Inefficient algorithms
#    - Unexpected library calls

# 4. In pprof web UI:
#    - Check "Flame Graph" for hot paths
#    - Look at "Top" sorted by "flat"
#    - Review "Source" for specific functions
```

### Scenario 2: Memory Leak

```bash
# 1. Capture heap profile
curl http://localhost:8080/debug/pprof/heap > heap1.prof

# 2. Wait and capture again
sleep 300
curl http://localhost:8080/debug/pprof/heap > heap2.prof

# 3. Compare
go tool pprof -http=:8080 -base=heap1.prof heap2.prof

# 4. Look for:
#    - Growing allocations
#    - Unexpected retention
#    - Caches not being cleared

# 5. Use inuse_space to see current memory
go tool pprof -sample_index=inuse_space -http=:8080 heap2.prof
```

### Scenario 3: Goroutine Leak

```bash
# 1. Check goroutine count
curl http://localhost:8080/api/stats

# 2. Capture goroutine profile
curl http://localhost:8080/debug/pprof/goroutine > goroutine.prof

# 3. Analyze
go tool pprof -http=:8080 goroutine.prof

# 4. Look for:
#    - Unexpected high count
#    - Blocked goroutines
#    - Waiting on channels

# 5. In text mode
go tool pprof -text goroutine.prof | grep -A 5 "chan receive"
```

### Scenario 4: Slow Response Times

```bash
# 1. Check for blocking
curl http://localhost:8080/debug/pprof/block > block.prof
go tool pprof -http=:8080 block.prof

# 2. Check for mutex contention
curl http://localhost:8080/debug/pprof/mutex > mutex.prof
go tool pprof -http=:8080 mutex.prof

# 3. Check goroutine scheduling with trace
curl http://localhost:8080/debug/pprof/trace?seconds=5 > trace.out
go tool trace trace.out

# 4. Look for:
#    - Lock contention
#    - Channel blocking
#    - Poor goroutine scheduling
```

### Scenario 5: GC Pressure

```bash
# 1. Capture allocation profile
curl http://localhost:8080/debug/pprof/allocs > allocs.prof

# 2. View total allocations
go tool pprof -sample_index=alloc_space -http=:8080 allocs.prof

# 3. Look for:
#    - Excessive allocations in hot paths
#    - String concatenation in loops
#    - Unnecessary conversions
#    - Not reusing buffers

# 4. Optimize by:
#    - Preallocating slices
#    - Using sync.Pool
#    - Reducing allocations in loops
```

### Scenario 6: Benchmark Optimization

```bash
# 1. Run benchmark with profiling
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof -benchmem

# 2. Analyze CPU
go tool pprof -http=:8080 cpu.prof

# 3. Analyze memory
go tool pprof -http=:8080 mem.prof

# 4. Make optimizations

# 5. Benchmark again
go test -bench=. -cpuprofile=cpu2.prof -memprofile=mem2.prof -benchmem

# 6. Compare
go tool pprof -base=cpu.prof cpu2.prof
go tool pprof -base=mem.prof mem2.prof
```

---

## Quick Reference

### Capture Profiles

| Profile Type | Web URL                           | CLI Flag                    |
| ------------ | --------------------------------- | --------------------------- |
| CPU          | `/debug/pprof/profile?seconds=30` | `-cpuprofile=cpu.prof`      |
| Heap         | `/debug/pprof/heap`               | `-memprofile=mem.prof`      |
| Goroutine    | `/debug/pprof/goroutine`          | `pprof.Lookup("goroutine")` |
| Block        | `/debug/pprof/block`              | `-blockprofile=block.prof`  |
| Mutex        | `/debug/pprof/mutex`              | `-mutexprofile=mutex.prof`  |
| Allocs       | `/debug/pprof/allocs`             | N/A                         |
| Thread       | `/debug/pprof/threadcreate`       | N/A                         |
| Trace        | N/A                               | `-trace=trace.out`          |

### Analysis Commands

| Command                                  | Description           |
| ---------------------------------------- | --------------------- |
| `go tool pprof -http=:8080 profile.prof` | Open web UI           |
| `go tool pprof profile.prof`             | Interactive CLI       |
| `top`                                    | Show top functions    |
| `top -cum`                               | Sort by cumulative    |
| `list funcName`                          | Show source           |
| `web`                                    | Open graph (graphviz) |
| `peek funcName`                          | Show callers/callees  |
| `traces`                                 | Show call traces      |

### Common Flags

| Flag                  | Description                 |
| --------------------- | --------------------------- |
| `-http=:8080`         | Start web UI                |
| `-base=old.prof`      | Compare against baseline    |
| `-focus=regex`        | Focus on matching functions |
| `-ignore=regex`       | Ignore matching functions   |
| `-sample_index=index` | Choose metric (heap)        |
| `-unit=unit`          | Display units (ms, MB, etc) |
| `-top`                | Show top entries            |
| `-text`               | Text output                 |

---

## Resources

- [Official pprof Documentation](https://pkg.go.dev/net/http/pprof)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [runtime/pprof Package](https://pkg.go.dev/runtime/pprof)
- [Execution Tracer](https://pkg.go.dev/runtime/trace)
- [Go Performance Workshop](https://github.com/davecheney/high-performance-go-workshop)

---

## Summary

1. **Always profile before optimizing** - Measure, don't guess
2. **Use the right profile type** - CPU, memory, goroutine, etc.
3. **Generate realistic load** - Profile under real conditions
4. **Use web UI** - `-http=:8080` for best experience
5. **Compare profiles** - Use `-base` to measure improvements
6. **Focus on hot paths** - Look at `flat` and `cum` times
7. **Combine profile types** - Get the complete picture
8. **Profile production safely** - Use authentication and localhost

Happy profiling! 🔥
