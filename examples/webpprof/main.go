package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	_ "net/http/pprof"
)

var (
	// Global state to demonstrate memory allocations
	userCache = make(map[int]*User)
	cacheMu   sync.RWMutex

	// Counter for operations
	requestCount uint64
	countMu      sync.Mutex
)

func main() {
	fmt.Println("Starting Web Application with pprof profiling...")
	fmt.Println("pprof endpoints available at http://localhost:8080/debug/pprof/")
	fmt.Println("")
	fmt.Println("Available endpoints:")
	fmt.Println("  http://localhost:8080/              - Home page")
	fmt.Println("  http://localhost:8080/api/users     - Create users (GET)")
	fmt.Println("  http://localhost:8080/api/compute   - CPU intensive task (GET)")
	fmt.Println("  http://localhost:8080/api/allocate  - Memory intensive task (GET)")
	fmt.Println("  http://localhost:8080/api/leak      - Simulate goroutine leak (GET)")
	fmt.Println("  http://localhost:8080/api/stats     - Application statistics (GET)")
	fmt.Println("")
	fmt.Println("pprof profiles:")
	fmt.Println("  http://localhost:8080/debug/pprof/              - Index")
	fmt.Println("  http://localhost:8080/debug/pprof/heap          - Heap profile")
	fmt.Println("  http://localhost:8080/debug/pprof/goroutine     - Goroutines")
	fmt.Println("  http://localhost:8080/debug/pprof/profile       - CPU profile (30s)")
	fmt.Println("  http://localhost:8080/debug/pprof/block         - Block profile")
	fmt.Println("  http://localhost:8080/debug/pprof/mutex         - Mutex profile")
	fmt.Println("  http://localhost:8080/debug/pprof/threadcreate  - Thread creation")
	fmt.Println("  http://localhost:8080/debug/pprof/allocs        - All memory allocations")

	// Enable profiling for blocking and mutex
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	// Setup routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/users", createUsersHandler)
	http.HandleFunc("/api/compute", computeHandler)
	http.HandleFunc("/api/allocate", allocateHandler)
	http.HandleFunc("/api/leak", goroutineLeakHandler)
	http.HandleFunc("/api/stats", statsHandler)

	// Start background workers
	go backgroundWorker()

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
