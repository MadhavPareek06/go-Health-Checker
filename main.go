package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const configFilePath = "config.json"

func main() {
	// Load configuration from file and environment variables.
	config, err := LoadConfiguration(configFilePath)
	if err!= nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Use a WaitGroup to ensure all goroutines finish before main exits.
	var wg sync.WaitGroup
	
	// Create channels for communication between goroutines.
	// jobs is a buffered channel to decouple the producer from the workers.
	jobs := make(chan string, config.Concurrency) 
	results := make(chan HealthCheckResult, config.Concurrency)

	// Create the components.
	checker := NewChecker(config.RequestTimeout)
	aggregator := NewAggregatedMetrics()
	reporter := &ConsoleReporter{}

	// Start the worker pool.
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			WorkerPool(checker, jobs, results)
		}()
	}

	// Graceful shutdown channel to receive OS signals.
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// A goroutine to listen for results and aggregate them.
	go func() {
		for result := range results {
			aggregator.AddResult(result)
		}
	}()

	// The main orchestration loop.
	go func() {
		ticker := time.NewTicker(config.PingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				for _, service := range config.Services {
					jobs <- service
				}
				// After sending all jobs, calculate and report the metrics.
				aggregator.CalculateAverageLatency()
				reporter.Report(aggregator)
			case <-shutdownChan:
				log.Println("Received shutdown signal. Stopping the health monitor.")
				close(jobs)
				return
			}
		}
	}()

	// Wait for the main orchestrator goroutine to finish.
	// The main goroutine will block here until it receives a shutdown signal.
	<-shutdownChan
	
	// Ensure all workers have finished processing remaining jobs.
	wg.Wait()
	
	close(results)
	log.Println("Health monitor gracefully shut down.")
}