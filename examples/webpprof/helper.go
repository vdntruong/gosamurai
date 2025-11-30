package main

import "time"

func fibonacciCompute(n int) uint64 {
	var result uint64
	for i := 0; i < n; i++ {
		result += fibonacci(20)
	}
	return result
}

func fibonacci(n int) uint64 {
	if n <= 1 {
		return uint64(n)
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func incrementCounter() {
	countMu.Lock()
	requestCount++
	countMu.Unlock()
}

func backgroundWorker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Simulate background work
		cacheMu.Lock()
		// Clean old entries if cache is too large
		if len(userCache) > 10000 {
			for id := range userCache {
				delete(userCache, id)
				break // Delete one at a time
			}
		}
		cacheMu.Unlock()
	}
}
