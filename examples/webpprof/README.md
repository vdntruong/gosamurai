# Web Application with pprof Profiling

A complete example demonstrating how to integrate pprof profiling into a web application.

## Features

This web application includes:
- **Multiple HTTP endpoints** with different workloads
- **Integrated pprof handlers** for all profile types
- **Sample workloads** for CPU, memory, and goroutines
- **Statistics endpoint** for runtime metrics

## Quick Start

```bash
# Run the application
go run main.go

# Application will start on http://localhost:8080
```

## Endpoints

### Application Endpoints

- `http://localhost:8080/` - Home page with links
- `http://localhost:8080/api/users?count=100` - Create users (memory allocation)
- `http://localhost:8080/api/compute?iterations=1000000` - CPU intensive task
- `http://localhost:8080/api/allocate?size=1000` - Allocate memory (MB)
- `http://localhost:8080/api/leak?count=10` - Simulate goroutine leak
- `http://localhost:8080/api/stats` - Runtime statistics

### pprof Endpoints

- `http://localhost:8080/debug/pprof/` - pprof index
- `http://localhost:8080/debug/pprof/heap` - Heap profile
- `http://localhost:8080/debug/pprof/goroutine` - Goroutine profile
- `http://localhost:8080/debug/pprof/profile?seconds=30` - CPU profile (30s)
- `http://localhost:8080/debug/pprof/block` - Block profile
- `http://localhost:8080/debug/pprof/mutex` - Mutex profile
- `http://localhost:8080/debug/pprof/allocs` - Allocation profile
- `http://localhost:8080/debug/pprof/threadcreate` - Thread creation profile

## Usage Examples

### 1. Generate Load

```bash
# Create users
curl "http://localhost:8080/api/users?count=1000"

# CPU intensive
curl "http://localhost:8080/api/compute?iterations=5000000"

# Allocate memory
curl "http://localhost:8080/api/allocate?size=500"

# Check stats
curl "http://localhost:8080/api/stats" | jq
```

### 2. Capture Profiles

```bash
# CPU profile (30 seconds)
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof

# Heap profile
curl http://localhost:8080/debug/pprof/heap > heap.prof

# Goroutine profile
curl http://localhost:8080/debug/pprof/goroutine > goroutine.prof

# Block profile
curl http://localhost:8080/debug/pprof/block > block.prof

# Mutex profile
curl http://localhost:8080/debug/pprof/mutex > mutex.prof
```

### 3. Analyze Profiles

```bash
# Web UI (recommended)
go tool pprof -http=:9090 cpu.prof

# Interactive CLI
go tool pprof cpu.prof
(pprof) top
(pprof) list main.computeHandler
(pprof) web

# Text output
go tool pprof -top cpu.prof
```

## Complete Workflow Example

```bash
# 1. Start the server
go run main.go &

# 2. Generate some load to warm up
for i in {1..10}; do
    curl "http://localhost:8080/api/users?count=100" &
    curl "http://localhost:8080/api/compute?iterations=100000" &
done
wait

# 3. Start CPU profiling and generate load
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof &
for i in {1..50}; do
    curl "http://localhost:8080/api/compute?iterations=1000000" &
done
wait

# 4. Capture memory profile
curl http://localhost:8080/debug/pprof/heap > heap.prof

# 5. Analyze
go tool pprof -http=:9090 cpu.prof
go tool pprof -http=:9091 heap.prof
```

## Load Testing

Use tools like `hey` or `ab` for better load testing:

```bash
# Install hey
go install github.com/rakyll/hey@latest

# Generate load while profiling
curl http://localhost:8080/debug/pprof/profile?seconds=30 > cpu.prof &
hey -n 10000 -c 100 http://localhost:8080/api/compute?iterations=50000
wait

# Analyze results
go tool pprof -http=:9090 cpu.prof
```

## What to Look For

### CPU Profile
- Functions with high `flat` time (exclusive time)
- Functions with high `cum` time (cumulative time)
- Recursive calls or tight loops

### Heap Profile
- Large allocations
- Functions allocating frequently
- Memory that's not being freed

### Goroutine Profile
- Growing number of goroutines
- Blocked goroutines
- Leaked goroutines

### Block Profile
- Channel operations blocking
- Lock contention
- Select statements

### Mutex Profile
- Mutex contention
- Lock wait times

## See Also

- [Complete pprof Guide](../PPROF_GUIDE.md) - Comprehensive documentation
- [CLI Example](../clipprof/) - CLI application with profiling
