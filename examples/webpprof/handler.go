package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<html>
		<head><title>pprof Web Example</title></head>
		<body>
			<h1>Web Application with pprof Profiling</h1>
			<h2>API Endpoints</h2>
			<ul>
				<li><a href="/api/users?count=100">Create 100 Users</a></li>
				<li><a href="/api/compute?iterations=1000000">CPU Intensive Task</a></li>
				<li><a href="/api/allocate?size=1000">Memory Allocation</a></li>
				<li><a href="/api/leak?count=10">Simulate Goroutine Leak</a></li>
				<li><a href="/api/stats">Application Statistics</a></li>
			</ul>
			<h2>pprof Profiles</h2>
			<ul>
				<li><a href="/debug/pprof/">pprof Index</a></li>
				<li><a href="/debug/pprof/heap">Heap Profile</a></li>
				<li><a href="/debug/pprof/goroutine">Goroutine Profile</a></li>
				<li><a href="/debug/pprof/profile?seconds=10">CPU Profile (10s)</a></li>
				<li><a href="/debug/pprof/allocs">Allocation Profile</a></li>
				<li><a href="/debug/pprof/block">Block Profile</a></li>
				<li><a href="/debug/pprof/mutex">Mutex Profile</a></li>
			</ul>
		</body>
		</html>
	`)
}

func createUsersHandler(w http.ResponseWriter, r *http.Request) {
	count := 100
	if c := r.URL.Query().Get("count"); c != "" {
		fmt.Sscanf(c, "%d", &count)
	}

	users := make([]*User, count)
	for i := 0; i < count; i++ {
		user := &User{
			ID:        i + 1,
			Name:      fmt.Sprintf("User %d", i+1),
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			CreatedAt: time.Now(),
			Metadata: map[string]interface{}{
				"role":   "user",
				"active": true,
				"score":  rand.Intn(100),
			},
		}
		users[i] = user

		// Store in cache
		cacheMu.Lock()
		userCache[user.ID] = user
		cacheMu.Unlock()
	}

	incrementCounter()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"count":   count,
		"message": fmt.Sprintf("Created %d users", count),
	})
}

func computeHandler(w http.ResponseWriter, r *http.Request) {
	iterations := 1000000
	if i := r.URL.Query().Get("iterations"); i != "" {
		fmt.Sscanf(i, "%d", &iterations)
	}

	start := time.Now()
	result := fibonacciCompute(iterations)
	duration := time.Since(start)

	incrementCounter()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"iterations": iterations,
		"result":     result,
		"duration":   duration.String(),
	})
}

func allocateHandler(w http.ResponseWriter, r *http.Request) {
	size := 1000
	if s := r.URL.Query().Get("size"); s != "" {
		fmt.Sscanf(s, "%d", &size)
	}

	// Allocate large slices to stress memory
	var data [][]byte
	for i := 0; i < size; i++ {
		chunk := make([]byte, 1024*1024) // 1MB per chunk
		for j := range chunk {
			chunk[j] = byte(rand.Intn(256))
		}
		data = append(data, chunk)
	}

	incrementCounter()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "success",
		"allocated_mb":   size,
		"heap_alloc_mb":  memStats.HeapAlloc / 1024 / 1024,
		"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
	})

	// Keep data alive until response is sent
	_ = data
}

func goroutineLeakHandler(w http.ResponseWriter, r *http.Request) {
	count := 10
	if c := r.URL.Query().Get("count"); c != "" {
		fmt.Sscanf(c, "%d", &count)
	}

	// Create goroutines that will never finish
	for i := 0; i < count; i++ {
		go func(id int) {
			ch := make(chan struct{})
			<-ch // Block forever
		}(i)
	}

	incrementCounter()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":            "success",
		"leaked_goroutines": count,
		"total_goroutines":  runtime.NumGoroutine(),
		"warning":           "These goroutines will leak!",
	})
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	cacheMu.RLock()
	cacheSize := len(userCache)
	cacheMu.RUnlock()

	countMu.Lock()
	count := requestCount
	countMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"goroutines":     runtime.NumGoroutine(),
		"heap_alloc_mb":  memStats.HeapAlloc / 1024 / 1024,
		"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
		"sys_mb":         memStats.Sys / 1024 / 1024,
		"gc_runs":        memStats.NumGC,
		"cache_size":     cacheSize,
		"request_count":  count,
	})
}
