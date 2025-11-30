package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync"
	"time"
)

var (
	cpuProfile   = flag.String("cpuprofile", "", "write cpu profile to file")
	memProfile   = flag.String("memprofile", "", "write memory profile to file")
	traceFile    = flag.String("trace", "", "write execution trace to file")
	blockProfile = flag.String("blockprofile", "", "write block profile to file")
	mutexProfile = flag.String("mutexprofile", "", "write mutex profile to file")

	workload   = flag.String("workload", "all", "workload type: cpu, memory, goroutines, all")
	iterations = flag.Int("iterations", 1000, "number of iterations for CPU workload")
	allocSize  = flag.Int("allocsize", 1000, "size in MB for memory workload")
	goroutines = flag.Int("goroutines", 100, "number of goroutines to spawn")
	duration   = flag.Int("duration", 10, "duration in seconds to run workload")
)

func main() {
	flag.Parse()

	fmt.Println("CLI Application with pprof Profiling")
	fmt.Println("=====================================")
	fmt.Printf("Workload: %s\n", *workload)
	fmt.Printf("Duration: %d seconds\n", *duration)
	fmt.Println()

	// Setup CPU profiling
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
		fmt.Printf("CPU profiling enabled, writing to: %s\n", *cpuProfile)
	}

	// Setup trace
	if *traceFile != "" {
		f, err := os.Create(*traceFile)
		if err != nil {
			log.Fatal("could not create trace file: ", err)
		}
		defer f.Close()

		if err := trace.Start(f); err != nil {
			log.Fatal("could not start trace: ", err)
		}
		defer trace.Stop()
		fmt.Printf("Execution trace enabled, writing to: %s\n", *traceFile)
	}

	// Enable block and mutex profiling
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	fmt.Println("\nStarting workload...")
	startTime := time.Now()

	// Run workload
	switch *workload {
	case "cpu":
		runCPUWorkload()
	case "memory":
		runMemoryWorkload()
	case "goroutines":
		runGoroutineWorkload()
	case "all":
		runAllWorkloads()
	default:
		log.Fatalf("Unknown workload: %s", *workload)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nWorkload completed in %s\n", elapsed)

	// Write memory profile
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()

		runtime.GC() // Get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		fmt.Printf("Memory profile written to: %s\n", *memProfile)
	}

	// Write block profile
	if *blockProfile != "" {
		f, err := os.Create(*blockProfile)
		if err != nil {
			log.Fatal("could not create block profile: ", err)
		}
		defer f.Close()

		if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
			log.Fatal("could not write block profile: ", err)
		}
		fmt.Printf("Block profile written to: %s\n", *blockProfile)
	}

	// Write mutex profile
	if *mutexProfile != "" {
		f, err := os.Create(*mutexProfile)
		if err != nil {
			log.Fatal("could not create mutex profile: ", err)
		}
		defer f.Close()

		if err := pprof.Lookup("mutex").WriteTo(f, 0); err != nil {
			log.Fatal("could not write mutex profile: ", err)
		}
		fmt.Printf("Mutex profile written to: %s\n", *mutexProfile)
	}

	// Print statistics
	printStats()
}

func runCPUWorkload() {
	fmt.Println("Running CPU-intensive workload...")
	endTime := time.Now().Add(time.Duration(*duration) * time.Second)

	var result uint64
	count := 0
	for time.Now().Before(endTime) {
		result += computeFibonacci(30)
		result += computePrimes(10000)
		count++
	}

	fmt.Printf("CPU workload: %d iterations, result: %d\n", count, result)
}

func runMemoryWorkload() {
	fmt.Println("Running memory-intensive workload...")

	// Allocate large chunks of memory
	var data [][]byte
	totalMB := 0

	for i := 0; i < *allocSize; i++ {
		chunk := make([]byte, 1024*1024) // 1MB
		// Fill with random data to prevent optimization
		for j := 0; j < len(chunk); j += 1024 {
			chunk[j] = byte(rand.Intn(256))
		}
		data = append(data, chunk)
		totalMB++

		if i%100 == 0 && i > 0 {
			fmt.Printf("Allocated %d MB...\n", totalMB)
		}
	}

	fmt.Printf("Memory workload: allocated %d MB\n", totalMB)

	// Keep data alive
	time.Sleep(time.Duration(*duration) * time.Second)
	_ = data
}

func runGoroutineWorkload() {
	fmt.Println("Running goroutine workload...")

	var wg sync.WaitGroup
	startChan := make(chan struct{})

	// Create goroutines
	for i := 0; i < *goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-startChan

			// Each goroutine does some work
			var result uint64
			endTime := time.Now().Add(time.Duration(*duration) * time.Second)
			for time.Now().Before(endTime) {
				result += computeFibonacci(20)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	fmt.Printf("Started %d goroutines\n", *goroutines)
	close(startChan) // Start all goroutines

	wg.Wait()
	fmt.Println("All goroutines completed")
}

func runAllWorkloads() {
	fmt.Println("Running all workloads concurrently...")

	var wg sync.WaitGroup

	// CPU workload
	wg.Add(1)
	go func() {
		defer wg.Done()
		runCPUWorkload()
	}()

	// Memory workload
	wg.Add(1)
	go func() {
		defer wg.Done()
		runMemoryWorkload()
	}()

	// Goroutine workload
	wg.Add(1)
	go func() {
		defer wg.Done()
		runGoroutineWorkload()
	}()

	wg.Wait()
}

// Helper functions

func computeFibonacci(n int) uint64 {
	if n <= 1 {
		return uint64(n)
	}
	return computeFibonacci(n-1) + computeFibonacci(n-2)
}

func computePrimes(max int) uint64 {
	var count uint64
	for i := 2; i < max; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func printStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Println("\n=== Runtime Statistics ===")
	fmt.Printf("Goroutines:        %d\n", runtime.NumGoroutine())
	fmt.Printf("Heap Allocated:    %d MB\n", memStats.HeapAlloc/1024/1024)
	fmt.Printf("Total Allocated:   %d MB\n", memStats.TotalAlloc/1024/1024)
	fmt.Printf("System Memory:     %d MB\n", memStats.Sys/1024/1024)
	fmt.Printf("GC Runs:           %d\n", memStats.NumGC)
	fmt.Printf("Last GC Time:      %s\n", time.Unix(0, int64(memStats.LastGC)))
}
