package checks

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
)

// Result represents a health check result
type Result struct {
	State     string  `json:"state"` // up, down, degraded
	LatencyMs int64   `json:"latencyMs,omitempty"`
	Error     string  `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Executor runs health checks
type Executor struct {
	httpClient *network.Client
}

// NewExecutor creates a new check executor
func NewExecutor(netClient *network.Client) *Executor {
	return &Executor{
		httpClient: netClient,
	}
}

// Execute runs a single check and returns the result
func (e *Executor) Execute(ctx context.Context, check config.Check) (*Result, error) {
	switch check.Type {
	case "http":
		return e.executeHTTP(ctx, check)
	case "tcp":
		return e.executeTCP(ctx, check)
	case "ping":
		return e.executePing(ctx, check)
	default:
		return &Result{
			State:     "down",
			Error:     fmt.Sprintf("unknown check type: %s", check.Type),
			Timestamp: time.Now(),
		}, nil
	}
}

// executeHTTP performs an HTTP health check
func (e *Executor) executeHTTP(ctx context.Context, check config.Check) (*Result, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", check.URL, nil)
	if err != nil {
		return &Result{
			State:     "down",
			Error:     err.Error(),
			Timestamp: time.Now(),
		}, nil
	}

	// Add custom headers
	for k, v := range check.Headers {
		req.Header.Set(k, v)
	}

	// Use custom client for SSRF protection
	resp, err := e.httpClient.Do(req)
	latency := time.Since(start)

	if err != nil {
		return &Result{
			State:     "down",
			Error:     err.Error(),
			LatencyMs: latency.Milliseconds(),
			Timestamp: time.Now(),
		}, nil
	}
	defer resp.Body.Close()

	// Check status code
	expectedStatus := check.ExpectStatus
	if expectedStatus == 0 {
		expectedStatus = 200
	}

	if resp.StatusCode != expectedStatus {
		return &Result{
			State:     "down",
			Error:     fmt.Sprintf("expected status %d, got %d", expectedStatus, resp.StatusCode),
			LatencyMs: latency.Milliseconds(),
			Timestamp: time.Now(),
		}, nil
	}

	// Check body content if specified
	if check.ExpectBodyContains != "" {
		// For MVP, skip body check (would need to read body)
		// TODO: Implement body content checking
	}

	// Determine state based on latency
	state := "up"
	if latency > 1*time.Second {
		state = "degraded"
	}

	return &Result{
		State:     state,
		LatencyMs: latency.Milliseconds(),
		Timestamp: time.Now(),
	}, nil
}

// executeTCP performs a TCP port check
func (e *Executor) executeTCP(ctx context.Context, check config.Check) (*Result, error) {
	if check.Host == "" {
		return &Result{
			State:     "down",
			Error:     "host is required",
			Timestamp: time.Now(),
		}, nil
	}
	if check.Port <= 0 || check.Port > 65535 {
		return &Result{
			State:     "down",
			Error:     fmt.Sprintf("invalid port: %d", check.Port),
			Timestamp: time.Now(),
		}, nil
	}

	start := time.Now()

	// Parse timeout
	timeout := 5 * time.Second
	if check.Timeout != "" {
		if d, err := time.ParseDuration(check.Timeout); err == nil {
			timeout = d
		}
	}

	dialer := &net.Dialer{
		Timeout: timeout,
	}

	address := fmt.Sprintf("%s:%d", check.Host, check.Port)
	conn, err := dialer.DialContext(ctx, "tcp", address)
	latency := time.Since(start)

	if err != nil {
		return &Result{
			State:     "down",
			Error:     err.Error(),
			LatencyMs: latency.Milliseconds(),
			Timestamp: time.Now(),
		}, nil
	}
	defer conn.Close()

	state := "up"
	if latency > 500*time.Millisecond {
		state = "degraded"
	}

	return &Result{
		State:     state,
		LatencyMs: latency.Milliseconds(),
		Timestamp: time.Now(),
	}, nil
}

// executePing performs a ping check (ICMP)
func (e *Executor) executePing(ctx context.Context, check config.Check) (*Result, error) {
	if check.Host == "" {
		return &Result{
			State:     "down",
			Error:     "host is required",
			Timestamp: time.Now(),
		}, nil
	}

	count := check.Count
	if count == 0 {
		count = 3
	}

	// Note: True ICMP ping requires raw sockets (elevated privileges)
	// For MVP, we'll use TCP ping to common ports as fallback
	// In production, use github.com/go-ping/ping with capabilities

	// Try TCP ping to port 80 or 443 as fallback
	ports := []int{80, 443, 22}
	var success int
	var totalLatency int64

	for _, port := range ports {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", check.Host, port), 2*time.Second)
		latency := time.Since(start)

		if err == nil {
			conn.Close()
			success++
			totalLatency += latency.Milliseconds()
			if success >= count {
				break
			}
		}
	}

	if success == 0 {
		return &Result{
			State:     "down",
			Error:     "host unreachable (ping requires elevated privileges or open ports)",
			Timestamp: time.Now(),
		}, nil
	}

	avgLatency := totalLatency / int64(success)

	state := "up"
	if avgLatency > 100 {
		state = "degraded"
	}

	return &Result{
		State:     state,
		LatencyMs: avgLatency,
		Timestamp: time.Now(),
	}, nil
}

// GetInterval returns the check interval in seconds
func GetInterval(check config.Check) int {
	if check.IntervalSeconds > 0 {
		return check.IntervalSeconds
	}
	return 60 // Default 60 seconds
}

// GetTimeout returns the check timeout duration
func GetTimeout(check config.Check) time.Duration {
	if check.Timeout != "" {
		if d, err := time.ParseDuration(check.Timeout); err == nil {
			return d
		}
	}
	return 10 * time.Second // Default 10 seconds
}

// GetRetries returns the number of retries
func GetRetries(check config.Check) int {
	if check.Retries > 0 {
		return check.Retries
	}
	return 1 // Default 1 retry
}

// StateFromError converts an error to a state string
func StateFromError(err error) string {
	if err == nil {
		return "up"
	}
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") {
		return "down"
	}
	return "degraded"
}
