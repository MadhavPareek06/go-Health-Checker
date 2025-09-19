package main

import (
	"fmt"
	"time"
)

// HealthReporter is an interface that defines the contract for reporting
// the health status of a service monitor. This allows for different
// implementations without changing the core logic.
type HealthReporter interface {
	Report(metrics *AggregatedMetrics)
}

// ConsoleReporter is a concrete implementation of HealthReporter
// that prints the results to the console.
type ConsoleReporter struct{}

// Report prints the aggregated metrics to the console.
func (r *ConsoleReporter) Report(metrics *AggregatedMetrics) {
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Health Monitor Report (%s)\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("--------------------------------------------------")
	
	fmt.Printf("Total Checks: %d\n", metrics.TotalChecks)
	fmt.Printf("Successful Pings: %d\n", metrics.SuccessfulPings)
	fmt.Printf("Failed Pings: %d\n", metrics.FailedPings)
	fmt.Printf("Average Latency: %s\n", metrics.AverageLatency)

	fmt.Println("--------------------------------------------------")
	fmt.Println("Service Status:")
	
	metrics.ServiceStatus.Range(func(key, value any) bool {
		url, _ := key.(string)
		status, _ := value.(ServiceMetrics)
		
		fmt.Printf("  - %s | Status: %s | Latency: %s\n", url, status.Status, status.Latency)
		return true
	})
	
	fmt.Println("--------------------------------------------------")
}