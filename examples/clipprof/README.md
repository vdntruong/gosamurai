# CLI Application with pprof Profiling

A complete example demonstrating how to add profiling to a CLI application using command-line flags.

## Features

- **Multiple workload types**: CPU, memory, goroutines, or all combined
- **Comprehensive profiling flags**: CPU, memory, block, mutex, trace
- **Configurable parameters**: iterations, size, goroutines, duration
- **Runtime statistics**: goroutines, memory usage, GC stats

## Quick Start

```bash
# Basic run with CPU profiling
go run main.go -cpuprofile=cpu.prof -workload=cpu -duration=10

# All workloads with multiple profiles
go run main.go \
  -cpuprofile=cpu.prof \
  -memprofile=mem.prof \
  -trace=trace.out \
  -workload=all \
  -duration=10
```

## Command-Line Flags

### Profiling Flags

- `-cpuprofile=<file>` - Enable CPU profiling, write to file
- `-memprofile=<file>` - Enable memory profiling, write to file
- `-blockprofile=<file>` - Enable block profiling, write to file
- `-mutexprofile=<file>` - Enable mutex profiling, write to file
- `-trace=<file>` - Enable execution trace, write to file

### Workload Flags

- `-workload=<type>` - Workload type: `cpu`, `memory`, `goroutines`, or `all` (default: `all`)
- `-iterations=<N>` - Iterations for CPU workload (default: 1000)
- `-allocsize=<MB>` - Size in MB for memory workload (default: 1000)
- `-goroutines=<N>` - Number of goroutines to spawn (default: 100)
- `-duration=<seconds>` - Duration in seconds to run workload (default: 10)

## Usage Examples

### CPU Profiling

```bash
# CPU intensive workload
go run main.go -cpuprofile=cpu.prof -workload=cpu -duration=10

# Analyze
go tool pprof -http=:8080 cpu.prof
```

### Memory Profiling

```bash
# Memory intensive workload
go run main.go -memprofile=mem.prof -workload=memory -allocsize=500

# Analyze in-use memory
go tool pprof -sample_index=inuse_space -http=:8080 mem.prof

# Analyze total allocations
go tool pprof -sample_index=alloc_space -http=:8080 mem.prof
```

### Goroutine Profiling

```bash
# Spawn many goroutines
go run main.go -workload=goroutines -goroutines=500 -duration=10

# To see goroutines, you need to write goroutine profile programmatically
# or use block/mutex profiles to see contention
go run main.go \
  -blockprofile=block.prof \
  -mutexprofile=mutex.prof \
  -workload=goroutines \
  -goroutines=500

# Analyze
go tool pprof -http=:8080 block.prof
```

### Multiple Profiles

```bash
# Capture all profile types
go run main.go \
  -cpuprofile=cpu.prof \
  -memprofile=mem.prof \
  -blockprofile=block.prof \
  -mutexprofile=mutex.prof \
  -trace=trace.out \
  -workload=all \
  -duration=15

# Analyze each
go tool pprof -http=:8080 cpu.prof
go tool pprof -http=:8081 mem.prof
go tool pprof -http=:8082 block.prof
go tool trace trace.out
```

### Execution Trace

```bash
# Capture execution trace
go run main.go -trace=trace.out -workload=all -duration=10

# View trace (opens browser)
go tool trace trace.out
```

## Complete Workflow Examples

### Example 1: CPU Optimization

```bash
# 1. Capture baseline
go run main.go -cpuprofile=cpu_before.prof -workload=cpu -iterations=10000

# 2. Make code changes

# 3. Capture after changes
go run main.go -cpuprofile=cpu_after.prof -workload=cpu -iterations=10000

# 4. Compare
go tool pprof -base=cpu_before.prof cpu_after.prof
(pprof) top -cum
```

### Example 2: Memory Analysis

```bash
# 1. Run with memory profiling
go run main.go \
  -memprofile=mem.prof \
  -workload=memory \
  -allocsize=1000 \
  -duration=10

# 2. Analyze current memory usage
go tool pprof -sample_index=inuse_space -http=:8080 mem.prof

# 3. Look at top allocators
go tool pprof -top -sample_index=alloc_space mem.prof
```

### Example 3: Goroutine Analysis

```bash
# 1. Run workload that creates goroutines
go run main.go \
  -blockprofile=block.prof \
  -mutexprofile=mutex.prof \
  -trace=trace.out \
  -workload=goroutines \
  -goroutines=1000 \
  -duration=10

# 2. Analyze blocking
go tool pprof -http=:8080 block.prof

# 3. Analyze mutex contention
go tool pprof -http=:8081 mutex.prof

# 4. View execution trace
go tool trace trace.out
```

### Example 4: Comparison Testing

```bash
# Baseline
go run main.go \
  -cpuprofile=v1_cpu.prof \
  -memprofile=v1_mem.prof \
  -workload=all \
  -duration=20

# After optimization
go run main.go \
  -cpuprofile=v2_cpu.prof \
  -memprofile=v2_mem.prof \
  -workload=all \
  -duration=20

# Compare CPU
go tool pprof -http=:8080 -base=v1_cpu.prof v2_cpu.prof

# Compare memory
go tool pprof -http=:8081 -base=v1_mem.prof v2_mem.prof
```

## Workload Types

### CPU Workload
- Computes Fibonacci numbers recursively
- Finds prime numbers
- Runs for specified duration

### Memory Workload
- Allocates large byte slices
- Fills them with random data
- Keeps data alive for duration

### Goroutines Workload
- Spawns specified number of goroutines
- Each does CPU work with sleep
- Uses channels for synchronization

### All Workload
- Runs all workloads concurrently
- Good for stress testing
- Captures diverse profile data

## Analyzing Results

### CPU Profile Analysis

```bash
# Web UI
go tool pprof -http=:8080 cpu.prof

# Command line
go tool pprof cpu.prof
(pprof) top          # Top functions by CPU
(pprof) top -cum     # Top by cumulative time
(pprof) list main.computeFibonacci  # Source code view
(pprof) web          # Graph visualization
```

### Memory Profile Analysis

```bash
# Current heap usage
go tool pprof -sample_index=inuse_space -http=:8080 mem.prof

# Total allocations
go tool pprof -sample_index=alloc_space -http=:8080 mem.prof

# Object counts
go tool pprof -sample_index=inuse_objects mem.prof
go tool pprof -sample_index=alloc_objects mem.prof
```

### Trace Analysis

```bash
# Open trace viewer (browser)
go tool trace trace.out

# In the viewer:
# - View trace: Timeline of all events
# - Goroutine analysis: Goroutine stats
# - Network blocking: Network wait times
# - Synchronization blocking: Lock contention
# - Syscall blocking: System call waits
# - Scheduler latency: Scheduling delays
```

## Tips

1. **Run long enough**: Short runs may not capture meaningful data
2. **Realistic workload**: Mirror production patterns when possible
3. **Multiple samples**: Run several times to average out noise
4. **Compare apples to apples**: Use same parameters when comparing
5. **Focus on your code**: Ignore runtime/stdlib in initial analysis

## Common Patterns

```bash
# Quick CPU check
go run main.go -cpuprofile=cpu.prof -workload=cpu -duration=5
go tool pprof -top cpu.prof

# Quick memory check
go run main.go -memprofile=mem.prof -workload=memory -allocsize=100
go tool pprof -top -sample_index=inuse_space mem.prof

# Full analysis
go run main.go \
  -cpuprofile=cpu.prof \
  -memprofile=mem.prof \
  -trace=trace.out \
  -workload=all \
  -duration=30
  
go tool pprof -http=:8080 cpu.prof &
go tool pprof -http=:8081 mem.prof &
go tool trace trace.out
```

## See Also

- [Complete pprof Guide](../PPROF_GUIDE.md) - Comprehensive documentation
- [Web Example](../webpprof/) - Web application with profiling
