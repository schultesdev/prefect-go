package errors

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "with message",
			err: &APIError{
				StatusCode: 404,
				Message:    "Flow not found",
			},
			want: "prefect api error (status 404): Flow not found",
		},
		{
			name: "without message",
			err: &APIError{
				StatusCode: 500,
			},
			want: "prefect api error (status 500)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Checks(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		checkFuncs map[string]func(*APIError) bool
		want       map[string]bool
	}{
		{
			name:       "404 not found",
			statusCode: 404,
			checkFuncs: map[string]func(*APIError) bool{
				"IsNotFound":     (*APIError).IsNotFound,
				"IsUnauthorized": (*APIError).IsUnauthorized,
				"IsForbidden":    (*APIError).IsForbidden,
				"IsRateLimited":  (*APIError).IsRateLimited,
				"IsServerError":  (*APIError).IsServerError,
			},
			want: map[string]bool{
				"IsNotFound":     true,
				"IsUnauthorized": false,
				"IsForbidden":    false,
				"IsRateLimited":  false,
				"IsServerError":  false,
			},
		},
		{
			name:       "401 unauthorized",
			statusCode: 401,
			checkFuncs: map[string]func(*APIError) bool{
				"IsNotFound":     (*APIError).IsNotFound,
				"IsUnauthorized": (*APIError).IsUnauthorized,
				"IsForbidden":    (*APIError).IsForbidden,
				"IsRateLimited":  (*APIError).IsRateLimited,
				"IsServerError":  (*APIError).IsServerError,
			},
			want: map[string]bool{
				"IsNotFound":     false,
				"IsUnauthorized": true,
				"IsForbidden":    false,
				"IsRateLimited":  false,
				"IsServerError":  false,
			},
		},
		{
			name:       "500 server error",
			statusCode: 500,
			checkFuncs: map[string]func(*APIError) bool{
				"IsNotFound":     (*APIError).IsNotFound,
				"IsUnauthorized": (*APIError).IsUnauthorized,
				"IsForbidden":    (*APIError).IsForbidden,
				"IsRateLimited":  (*APIError).IsRateLimited,
				"IsServerError":  (*APIError).IsServerError,
			},
			want: map[string]bool{
				"IsNotFound":     false,
				"IsUnauthorized": false,
				"IsForbidden":    false,
				"IsRateLimited":  false,
				"IsServerError":  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode}
			for checkName, checkFunc := range tt.checkFuncs {
				if got := checkFunc(err); got != tt.want[checkName] {
					t.Errorf("%s() = %v, want %v", checkName, got, tt.want[checkName])
				}
			}
		})
	}
}

func TestNewAPIError(t *testing.T) {
	tests := []struct {
		name       string
		resp       *http.Response
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "nil response",
			resp:       nil,
			wantStatus: 0,
			wantMsg:    "nil response",
		},
		{
			name: "404 with JSON body",
			resp: &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader(`{"detail": "Flow not found"}`)),
				Header:     http.Header{},
			},
			wantStatus: 404,
			wantMsg:    "Flow not found",
		},
		{
			name: "500 without body",
			resp: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     http.Header{},
			},
			wantStatus: 500,
			wantMsg:    "Internal Server Error",
		},
		{
			name: "with request ID",
			resp: &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(strings.NewReader(`{"detail": "Bad request"}`)),
				Header: http.Header{
					"X-Request-ID": []string{"req-123"},
				},
			},
			wantStatus: 400,
			wantMsg:    "Bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.resp)
			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("expected *APIError, got %T", err)
			}
			if apiErr.StatusCode != tt.wantStatus {
				t.Errorf("StatusCode = %v, want %v", apiErr.StatusCode, tt.wantStatus)
			}
			if apiErr.Message != tt.wantMsg {
				t.Errorf("Message = %v, want %v", apiErr.Message, tt.wantMsg)
			}
		})
	}
}

func TestNewAPIError_Validation(t *testing.T) {
	resp := &http.Response{
		StatusCode: 422,
		Body: io.NopCloser(strings.NewReader(`{
			"detail": [
				{
					"loc": ["body", "name"],
					"msg": "field required",
					"type": "value_error.missing"
				}
			]
		}`)),
		Header: http.Header{},
	}

	err := NewAPIError(resp)
	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}

	if valErr.StatusCode != 422 {
		t.Errorf("StatusCode = %v, want 422", valErr.StatusCode)
	}

	if len(valErr.ValidationErrors) != 1 {
		t.Fatalf("expected 1 validation error, got %d", len(valErr.ValidationErrors))
	}

	vd := valErr.ValidationErrors[0]
	if len(vd.Loc) != 2 || vd.Loc[0] != "body" || vd.Loc[1] != "name" {
		t.Errorf("Loc = %v, want [body name]", vd.Loc)
	}
	if vd.Msg != "field required" {
		t.Errorf("Msg = %v, want 'field required'", vd.Msg)
	}
	if vd.Type != "value_error.missing" {
		t.Errorf("Type = %v, want 'value_error.missing'", vd.Type)
	}
}

func TestHelperFunctions(t *testing.T) {
	notFoundErr := &APIError{StatusCode: 404}
	unauthorizedErr := &APIError{StatusCode: 401}
	serverErr := &APIError{StatusCode: 500}
	valErr := &ValidationError{APIError: &APIError{StatusCode: 422}}

	if !IsNotFound(notFoundErr) {
		t.Error("IsNotFound should return true for 404 error")
	}
	if !IsUnauthorized(unauthorizedErr) {
		t.Error("IsUnauthorized should return true for 401 error")
	}
	if !IsServerError(serverErr) {
		t.Error("IsServerError should return true for 500 error")
	}
	if !IsValidationError(valErr) {
		t.Error("IsValidationError should return true for validation error")
	}
}
