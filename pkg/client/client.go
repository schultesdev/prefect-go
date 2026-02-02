// Package client provides a Go client for the Prefect Server API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gregor/prefect-go/pkg/errors"
	"github.com/gregor/prefect-go/pkg/retry"
)

const (
	defaultBaseURL = "http://localhost:4200/api"
	defaultTimeout = 30 * time.Second
	userAgent      = "prefect-go-client/1.0.0"
)

// Client is the main client for interacting with the Prefect API.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiKey     string
	retrier    *retry.Retrier
	userAgent  string

	// Services
	Flows       *FlowsService
	FlowRuns    *FlowRunsService
	Deployments *DeploymentsService
	TaskRuns    *TaskRunsService
	WorkPools   *WorkPoolsService
	WorkQueues  *WorkQueuesService
	Variables   *VariablesService
	Logs        *LogsService
	Admin       *AdminService
}

// Option is a functional option for configuring the client.
type Option func(*Client) error

// NewClient creates a new Prefect API client with the given options.
func NewClient(opts ...Option) (*Client, error) {
	// Parse default base URL
	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default base URL: %w", err)
	}

	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		retrier:   retry.New(retry.DefaultConfig()),
		userAgent: userAgent,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// Initialize services
	c.Flows = &FlowsService{client: c}
	c.FlowRuns = &FlowRunsService{client: c}
	c.Deployments = &DeploymentsService{client: c}
	c.TaskRuns = &TaskRunsService{client: c}
	c.WorkPools = &WorkPoolsService{client: c}
	c.WorkQueues = &WorkQueuesService{client: c}
	c.Variables = &VariablesService{client: c}
	c.Logs = &LogsService{client: c}
	c.Admin = &AdminService{client: c}

	return c, nil
}

// WithBaseURL sets the base URL for the Prefect API.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("invalid base URL: %w", err)
		}
		c.baseURL = u
		return nil
	}
}

// WithAPIKey sets the API key for authentication with Prefect Cloud.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) error {
		c.apiKey = apiKey
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("HTTP client cannot be nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}

// WithRetry sets the retry configuration.
func WithRetry(config retry.Config) Option {
	return func(c *Client) error {
		c.retrier = retry.New(config)
		return nil
	}
}

// WithUserAgent sets a custom user agent string.
func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		c.userAgent = ua
		return nil
	}
}

// do executes an HTTP request with retry logic.
func (c *Client) do(ctx context.Context, method, path string, body, result interface{}) error {
	// Build full URL
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return fmt.Errorf("failed to parse path: %w", err)
	}

	// Serialize request body if present
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	// Set API key if present
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	// Execute request with retry logic
	var resp *http.Response
	resp, err = c.retrier.Do(ctx, func() (*http.Response, error) {
		// Clone the request for retry attempts
		reqClone := req.Clone(ctx)
		if reqBody != nil {
			// Reset body for retries
			if seeker, ok := reqBody.(io.Seeker); ok {
				if _, err := seeker.Seek(0, 0); err != nil {
					return nil, fmt.Errorf("failed to reset request body: %w", err)
				}
			}
			reqClone.Body = io.NopCloser(reqBody)
		}
		return c.httpClient.Do(reqClone)
	})

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return errors.NewAPIError(resp)
	}

	// Deserialize response body if result is provided
	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string, result interface{}) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, body, result interface{}) error {
	return c.do(ctx, http.MethodPost, path, body, result)
}

// patch performs a PATCH request.
func (c *Client) patch(ctx context.Context, path string, body, result interface{}) error {
	return c.do(ctx, http.MethodPatch, path, body, result)
}

// delete performs a DELETE request.
func (c *Client) delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// buildPath builds a URL path with optional query parameters.
func buildPath(base string, params map[string]string) string {
	if len(params) == 0 {
		return base
	}

	query := url.Values{}
	for k, v := range params {
		if v != "" {
			query.Set(k, v)
		}
	}

	if queryStr := query.Encode(); queryStr != "" {
		return base + "?" + queryStr
	}

	return base
}

// buildPathWithValues builds a URL path with url.Values query parameters.
func buildPathWithValues(base string, query url.Values) string {
	if len(query) == 0 {
		return base
	}

	if queryStr := query.Encode(); queryStr != "" {
		return base + "?" + queryStr
	}

	return base
}

// joinPath joins URL path segments.
func joinPath(parts ...string) string {
	return strings.Join(parts, "/")
}
