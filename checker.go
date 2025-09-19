package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HealthCheckResult represents the outcome of a single health check.
type HealthCheckResult struct {
	URL        string
	StatusCode int
	Latency    time.Duration
	Error      error
}

// Checker is a type that can perform health checks on a web service.
// It holds a reference to an http.Client.
type Checker struct {
	client *http.Client
}

// NewChecker creates a new Checker instance with a configured HTTP client.
// The client is configured with a timeout to prevent hanging requests.
func NewChecker(timeout time.Duration) *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Ping performs a health check on the given URL.
// It returns a HealthCheckResult with the outcome.
// The method uses a Context with a Timeout to ensure the request does not hang.
func (c *Checker) Ping(ctx context.Context, url string) HealthCheckResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err!= nil {
		return HealthCheckResult{
			URL:   url,
			Error: fmt.Errorf("failed to create HTTP request: %w", err),
		}
	}

	resp, err := c.client.Do(req)
	latency := time.Since(start)

	if err!= nil {
		return HealthCheckResult{
			URL:   url,
			Error: fmt.Errorf("failed to ping URL: %w", err),
		}
	}

	defer resp.Body.Close()

	return HealthCheckResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		Latency:    latency,
		Error:      nil,
	}
}

// WorkerPool starts a fixed number of goroutines that read from a jobs channel,
// perform health checks, and send results to a results channel.
func WorkerPool(checker *Checker, jobs <-chan string, results chan<- HealthCheckResult) {
	for url := range jobs {
		// Use a context for each ping, even though the client has a timeout,
		// to allow for future context-based cancellation or deadlines.
		ctx, cancel := context.WithTimeout(context.Background(), checker.client.Timeout)
		defer cancel()

		result := checker.Ping(ctx, url)
		results <- result
	}
}