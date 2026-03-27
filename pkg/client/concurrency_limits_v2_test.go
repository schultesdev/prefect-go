package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

func TestConcurrencyLimitsV2Service_Create(t *testing.T) {
	expected := models.ConcurrencyLimitV2{
		ID:    uuid.New(),
		Name:  "test-limit",
		Limit: 5,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/concurrency_limits/" {
			t.Errorf("path = %v, want /api/v2/concurrency_limits/", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	limit, err := client.ConcurrencyLimitsV2.Create(ctx, &models.ConcurrencyLimitV2Create{
		Name:  "test-limit",
		Limit: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limit.Name != "test-limit" {
		t.Errorf("Name = %v, want test-limit", limit.Name)
	}
}

func TestConcurrencyLimitsV2Service_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %v, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ConcurrencyLimitV2{ID: uuid.New(), Name: "test"})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	limit, err := client.ConcurrencyLimitsV2.Get(ctx, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limit.Name != "test" {
		t.Errorf("Name = %v, want test", limit.Name)
	}
}

func TestConcurrencyLimitsV2Service_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	err := client.ConcurrencyLimitsV2.Delete(ctx, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConcurrencyLimitsV2Service_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.GlobalConcurrencyLimitResponse{
			{ID: uuid.New(), Name: "test", Limit: 5, ActiveSlots: 0},
		})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	limits, err := client.ConcurrencyLimitsV2.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(limits) != 1 {
		t.Errorf("count = %v, want 1", len(limits))
	}
}

func TestConcurrencyLimitsV2Service_RenewLease(t *testing.T) {
	leaseID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	err := client.ConcurrencyLimitsV2.RenewLease(ctx, leaseID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
