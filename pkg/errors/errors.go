// Package errors provides structured error types for the Prefect API client.
package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error returned by the Prefect API.
type APIError struct {
	// StatusCode is the HTTP status code
	StatusCode int

	// Message is the error message
	Message string

	// Details contains additional error details from the API
	Details map[string]interface{}

	// RequestID is the request ID if available
	RequestID string

	// Response is the raw HTTP response
	Response *http.Response
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("prefect api error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("prefect api error (status %d)", e.StatusCode)
}

// IsNotFound returns true if the error is a 404 Not Found error.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden error.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsRateLimited returns true if the error is a 429 Too Many Requests error.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsServerError returns true if the error is a 5xx server error.
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// ValidationError represents a validation error from the API.
type ValidationError struct {
	*APIError
	ValidationErrors []ValidationDetail
}

// ValidationDetail represents a single validation error detail.
type ValidationDetail struct {
	Loc  []string    `json:"loc"`
	Msg  string      `json:"msg"`
	Type string      `json:"type"`
	Ctx  interface{} `json:"ctx,omitempty"`
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if len(e.ValidationErrors) == 0 {
		return e.APIError.Error()
	}
	return fmt.Sprintf("validation error: %s", e.ValidationErrors[0].Msg)
}

// NewAPIError creates a new APIError from an HTTP response.
func NewAPIError(resp *http.Response) error {
	if resp == nil {
		return &APIError{
			StatusCode: 0,
			Message:    "nil response",
		}
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Response:   resp,
		RequestID:  resp.Header.Get("X-Request-ID"),
	}

	// Try to read and parse the response body
	if resp.Body != nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil && len(body) > 0 {
			// Try to parse as JSON
			var errorResponse struct {
				Detail interface{} `json:"detail"`
			}
			if json.Unmarshal(body, &errorResponse) == nil {
				// Check if detail is a validation error (422)
				if resp.StatusCode == http.StatusUnprocessableEntity {
					return parseValidationError(apiErr, errorResponse.Detail)
				}

				// Handle string detail
				if msg, ok := errorResponse.Detail.(string); ok {
					apiErr.Message = msg
				} else if details, ok := errorResponse.Detail.(map[string]interface{}); ok {
					apiErr.Details = details
					if msg, ok := details["message"].(string); ok {
						apiErr.Message = msg
					}
				}
			} else {
				// Not JSON, use raw body as message
				apiErr.Message = string(body)
			}
		}
	}

	// Fallback to status text if no message
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}

	return apiErr
}

// parseValidationError parses validation errors from the API response.
func parseValidationError(apiErr *APIError, detail interface{}) error {
	valErr := &ValidationError{
		APIError: apiErr,
	}

	// Detail can be a list of validation errors
	if details, ok := detail.([]interface{}); ok {
		for _, d := range details {
			if detailMap, ok := d.(map[string]interface{}); ok {
				var vd ValidationDetail

				// Parse loc
				if loc, ok := detailMap["loc"].([]interface{}); ok {
					for _, l := range loc {
						if s, ok := l.(string); ok {
							vd.Loc = append(vd.Loc, s)
						}
					}
				}

				// Parse msg
				if msg, ok := detailMap["msg"].(string); ok {
					vd.Msg = msg
				}

				// Parse type
				if typ, ok := detailMap["type"].(string); ok {
					vd.Type = typ
				}

				// Parse ctx
				if ctx, ok := detailMap["ctx"]; ok {
					vd.Ctx = ctx
				}

				valErr.ValidationErrors = append(valErr.ValidationErrors, vd)
			}
		}
	}

	return valErr
}

// IsNotFound checks if an error is a 404 Not Found error.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}

// IsUnauthorized checks if an error is a 401 Unauthorized error.
func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.IsUnauthorized()
	}
	return false
}

// IsForbidden checks if an error is a 403 Forbidden error.
func IsForbidden(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.IsForbidden()
	}
	return false
}

// IsRateLimited checks if an error is a 429 Too Many Requests error.
func IsRateLimited(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.IsRateLimited()
	}
	return false
}

// IsServerError checks if an error is a 5xx server error.
func IsServerError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.IsServerError()
	}
	return false
}

// IsValidationError checks if an error is a validation error.
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
