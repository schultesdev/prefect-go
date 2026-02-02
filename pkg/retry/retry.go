// Package retry provides retry logic with exponential backoff for HTTP requests.
package retry

import (
	"context"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Config defines the configuration for the retry mechanism.
type Config struct {
	// MaxAttempts is the maximum number of retry attempts (including the first try).
	// Default: 3
	MaxAttempts int

	// InitialDelay is the initial delay between retries.
	// Default: 1 second
	InitialDelay time.Duration

	// MaxDelay is the maximum delay between retries.
	// Default: 30 seconds
	MaxDelay time.Duration

	// BackoffFactor is the multiplier for exponential backoff.
	// Default: 2.0 (exponential)
	BackoffFactor float64

	// RetryableErrors are HTTP status codes that should trigger a retry.
	// Default: 429, 500, 502, 503, 504
	RetryableErrors []int
}

// DefaultConfig returns the default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:     3,
		InitialDelay:    time.Second,
		MaxDelay:        30 * time.Second,
		BackoffFactor:   2.0,
		RetryableErrors: []int{429, 500, 502, 503, 504},
	}
}

// Retrier handles retry logic for HTTP requests.
type Retrier struct {
	config Config
}

// New creates a new Retrier with the given configuration.
func New(config Config) *Retrier {
	// Apply defaults if not set
	if config.MaxAttempts == 0 {
		config.MaxAttempts = DefaultConfig().MaxAttempts
	}
	if config.InitialDelay == 0 {
		config.InitialDelay = DefaultConfig().InitialDelay
	}
	if config.MaxDelay == 0 {
		config.MaxDelay = DefaultConfig().MaxDelay
	}
	if config.BackoffFactor == 0 {
		config.BackoffFactor = DefaultConfig().BackoffFactor
	}
	if len(config.RetryableErrors) == 0 {
		config.RetryableErrors = DefaultConfig().RetryableErrors
	}

	return &Retrier{
		config: config,
	}
}

// Do executes the given function with retry logic.
// It returns the response and any error encountered.
func (r *Retrier) Do(ctx context.Context, fn func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Execute the function
		resp, err = fn()

		// If no error and response is not retryable, return success
		if err == nil && !r.isRetryable(resp) {
			return resp, nil
		}

		// If this was the last attempt, return the error
		if attempt == r.config.MaxAttempts {
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		// Calculate delay with exponential backoff and jitter
		delay := r.calculateDelay(attempt, resp)

		// Wait before retrying
		select {
		case <-time.After(delay):
			// Continue to next attempt
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return resp, err
}

// isRetryable checks if the HTTP response indicates a retryable error.
func (r *Retrier) isRetryable(resp *http.Response) bool {
	if resp == nil {
		return false
	}

	for _, code := range r.config.RetryableErrors {
		if resp.StatusCode == code {
			return true
		}
	}

	return false
}

// calculateDelay calculates the delay before the next retry attempt.
// It implements exponential backoff with jitter and respects Retry-After header.
func (r *Retrier) calculateDelay(attempt int, resp *http.Response) time.Duration {
	// Check for Retry-After header (for 429 rate limiting)
	if resp != nil && resp.StatusCode == 429 {
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			// Try parsing as seconds
			if seconds, err := strconv.ParseInt(retryAfter, 10, 64); err == nil {
				delay := time.Duration(seconds) * time.Second
				if delay > r.config.MaxDelay {
					delay = r.config.MaxDelay
				}
				return delay
			}
			// Try parsing as HTTP date
			if t, err := http.ParseTime(retryAfter); err == nil {
				delay := time.Until(t)
				if delay < 0 {
					delay = 0
				}
				if delay > r.config.MaxDelay {
					delay = r.config.MaxDelay
				}
				return delay
			}
		}
	}

	// Calculate exponential backoff
	backoff := float64(r.config.InitialDelay) * math.Pow(r.config.BackoffFactor, float64(attempt-1))

	// Apply maximum delay cap
	if backoff > float64(r.config.MaxDelay) {
		backoff = float64(r.config.MaxDelay)
	}

	// Add jitter (random factor between 0.5 and 1.0)
	jitter := 0.5 + rand.Float64()*0.5
	delay := time.Duration(backoff * jitter)

	return delay
}
