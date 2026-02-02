package retry

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   Config
	}{
		{
			name:   "default config",
			config: Config{},
			want:   DefaultConfig(),
		},
		{
			name: "custom config",
			config: Config{
				MaxAttempts:     5,
				InitialDelay:    2 * time.Second,
				MaxDelay:        60 * time.Second,
				BackoffFactor:   3.0,
				RetryableErrors: []int{500},
			},
			want: Config{
				MaxAttempts:     5,
				InitialDelay:    2 * time.Second,
				MaxDelay:        60 * time.Second,
				BackoffFactor:   3.0,
				RetryableErrors: []int{500},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.config)
			if r.config.MaxAttempts != tt.want.MaxAttempts {
				t.Errorf("MaxAttempts = %v, want %v", r.config.MaxAttempts, tt.want.MaxAttempts)
			}
			if r.config.InitialDelay != tt.want.InitialDelay {
				t.Errorf("InitialDelay = %v, want %v", r.config.InitialDelay, tt.want.InitialDelay)
			}
		})
	}
}

func TestRetrier_isRetryable(t *testing.T) {
	r := New(DefaultConfig())

	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"429 rate limit", 429, true},
		{"500 server error", 500, true},
		{"502 bad gateway", 502, true},
		{"503 service unavailable", 503, true},
		{"504 gateway timeout", 504, true},
		{"200 ok", 200, false},
		{"400 bad request", 400, false},
		{"404 not found", 404, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			if got := r.isRetryable(resp); got != tt.want {
				t.Errorf("isRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetrier_Do_Success(t *testing.T) {
	r := New(DefaultConfig())
	ctx := context.Background()

	attempts := 0
	fn := func() (*http.Response, error) {
		attempts++
		return &http.Response{StatusCode: 200}, nil
	}

	resp, err := r.Do(ctx, fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1", attempts)
	}
}

func TestRetrier_Do_Retry(t *testing.T) {
	config := Config{
		MaxAttempts:     3,
		InitialDelay:    10 * time.Millisecond,
		MaxDelay:        100 * time.Millisecond,
		BackoffFactor:   2.0,
		RetryableErrors: []int{500},
	}
	r := New(config)
	ctx := context.Background()

	attempts := 0
	fn := func() (*http.Response, error) {
		attempts++
		if attempts < 3 {
			return &http.Response{StatusCode: 500}, nil
		}
		return &http.Response{StatusCode: 200}, nil
	}

	start := time.Now()
	resp, err := r.Do(ctx, fn)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
	// Should have waited at least once
	if elapsed < 10*time.Millisecond {
		t.Errorf("elapsed time %v too short, expected at least 10ms", elapsed)
	}
}

func TestRetrier_Do_ContextCancellation(t *testing.T) {
	config := Config{
		MaxAttempts:     5,
		InitialDelay:    100 * time.Millisecond,
		MaxDelay:        1 * time.Second,
		BackoffFactor:   2.0,
		RetryableErrors: []int{500},
	}
	r := New(config)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	attempts := 0
	fn := func() (*http.Response, error) {
		attempts++
		return &http.Response{StatusCode: 500}, nil
	}

	_, err := r.Do(ctx, fn)
	if err != context.DeadlineExceeded {
		t.Errorf("error = %v, want context.DeadlineExceeded", err)
	}
	if attempts > 2 {
		t.Errorf("attempts = %d, should be cancelled early", attempts)
	}
}

func TestRetrier_calculateDelay(t *testing.T) {
	config := Config{
		InitialDelay:  time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
	}
	r := New(config)

	tests := []struct {
		name    string
		attempt int
		resp    *http.Response
		minWant time.Duration
		maxWant time.Duration
	}{
		{
			name:    "first retry",
			attempt: 1,
			resp:    &http.Response{StatusCode: 500},
			minWant: 500 * time.Millisecond, // with jitter min 0.5
			maxWant: 1 * time.Second,        // with jitter max 1.0
		},
		{
			name:    "second retry",
			attempt: 2,
			resp:    &http.Response{StatusCode: 500},
			minWant: 1 * time.Second, // 2^1 * 1s * 0.5
			maxWant: 2 * time.Second, // 2^1 * 1s * 1.0
		},
		{
			name:    "with Retry-After header",
			attempt: 1,
			resp: &http.Response{
				StatusCode: 429,
				Header:     http.Header{"Retry-After": []string{"5"}},
			},
			minWant: 5 * time.Second,
			maxWant: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := r.calculateDelay(tt.attempt, tt.resp)
			if delay < tt.minWant || delay > tt.maxWant {
				t.Errorf("delay = %v, want between %v and %v", delay, tt.minWant, tt.maxWant)
			}
		})
	}
}
