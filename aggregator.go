package main

import (
	"sync"
	"sync/atomic"
	"time"
)

// ServiceMetrics holds the latest health status and latency for a single service.
type ServiceMetrics struct {
	Status  string
	Latency time.Duration
}

// AggregatedMetrics holds the summary of all health checks.
// It uses atomic counters and a sync.Map for thread-safe access.
type AggregatedMetrics struct {
	TotalChecks     int64
	SuccessfulPings int64
	FailedPings     int64
	TotalLatency    int64 // Stored in nanoseconds for atomic operations.
	AverageLatency  time.Duration
	ServiceStatus   sync.Map
}

// NewAggregatedMetrics creates a new, zero-valued AggregatedMetrics struct.
func NewAggregatedMetrics() *AggregatedMetrics {
	return &AggregatedMetrics{}
}

// AddResult aggregates a single health check result into the metrics.
// It uses atomic operations for counters and a sync.Map for service status.
func (a *AggregatedMetrics) AddResult(result HealthCheckResult) {
	// Increment total checks using an atomic counter.
	atomic.AddInt64(&a.TotalChecks, 1)

	if result.Error!= nil {
		atomic.AddInt64(&a.FailedPings, 1)
		a.ServiceStatus.Store(result.URL, ServiceMetrics{
			Status: "DOWN",
			Latency: 0,
		})
	} else {
		atomic.AddInt64(&a.SuccessfulPings, 1)
		atomic.AddInt64(&a.TotalLatency, result.Latency.Nanoseconds())
		a.ServiceStatus.Store(result.URL, ServiceMetrics{
			Status:  "UP",
			Latency: result.Latency,
		})
	}
}

// CalculateAverageLatency computes the average latency for successful pings.
func (a *AggregatedMetrics) CalculateAverageLatency() {
	successfulPings := atomic.LoadInt64(&a.SuccessfulPings)
	if successfulPings > 0 {
		totalLatency := atomic.LoadInt64(&a.TotalLatency)
		averageNs := totalLatency / successfulPings
		a.AverageLatency = time.Duration(averageNs) * time.Nanosecond
	} else {
		a.AverageLatency = 0
	}
}