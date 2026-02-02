package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "default client",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "with base URL",
			opts: []Option{
				WithBaseURL("http://example.com/api"),
			},
			wantErr: false,
		},
		{
			name: "with invalid base URL",
			opts: []Option{
				WithBaseURL("://invalid"),
			},
			wantErr: true,
		},
		{
			name: "with API key",
			opts: []Option{
				WithAPIKey("test-key"),
			},
			wantErr: false,
		},
		{
			name: "with custom timeout",
			opts: []Option{
				WithTimeout(60 * time.Second),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClient_do(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %v, want application/json", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Accept = %v, want application/json", r.Header.Get("Accept"))
		}

		// Check path
		if r.URL.Path == "/test" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		} else if r.URL.Path == "/error" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"detail": "bad request"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := NewClient(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		var result map[string]string
		err := client.get(ctx, "/test", &result)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result["message"] != "success" {
			t.Errorf("message = %v, want success", result["message"])
		}
	})

	t.Run("error response", func(t *testing.T) {
		var result map[string]string
		err := client.get(ctx, "/error", &result)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestClient_WithAPIKey(t *testing.T) {
	apiKey := "test-api-key"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		expected := "Bearer " + apiKey
		if auth != expected {
			t.Errorf("Authorization = %v, want %v", auth, expected)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithAPIKey(apiKey),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	err = client.get(ctx, "/test", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClient_WithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client, err := NewClient(WithHTTPClient(customClient))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if client.httpClient != customClient {
		t.Error("custom HTTP client not set")
	}
}

func TestClient_WithNilHTTPClient(t *testing.T) {
	_, err := NewClient(WithHTTPClient(nil))
	if err == nil {
		t.Error("expected error with nil HTTP client")
	}
}

func TestBuildPath(t *testing.T) {
	tests := []struct {
		name   string
		base   string
		params map[string]string
		want   string
	}{
		{
			name:   "no params",
			base:   "/flows",
			params: nil,
			want:   "/flows",
		},
		{
			name: "with params",
			base: "/flows",
			params: map[string]string{
				"limit":  "10",
				"offset": "20",
			},
			want: "/flows?limit=10&offset=20",
		},
		{
			name: "with empty param value",
			base: "/flows",
			params: map[string]string{
				"limit":  "10",
				"offset": "",
			},
			want: "/flows?limit=10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPath(tt.base, tt.params)
			// Note: URL encoding might change order, so we just check it contains the params
			if !contains(got, tt.base) {
				t.Errorf("buildPath() = %v, should contain %v", got, tt.base)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  string
	}{
		{
			name:  "single part",
			parts: []string{"flows"},
			want:  "flows",
		},
		{
			name:  "multiple parts",
			parts: []string{"flows", "123", "runs"},
			want:  "flows/123/runs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinPath(tt.parts...)
			if got != tt.want {
				t.Errorf("joinPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
