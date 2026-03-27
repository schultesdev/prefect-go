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

func TestConcurrencyLimitsService_Create(t *testing.T) {
	expected := models.ConcurrencyLimit{
		ID:               uuid.New(),
		Tag:              "test-tag",
		ConcurrencyLimit: 10,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	limit, err := client.ConcurrencyLimits.Create(ctx, &models.ConcurrencyLimitCreate{
		Tag:              "test-tag",
		ConcurrencyLimit: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limit.Tag != "test-tag" {
		t.Errorf("Tag = %v, want test-tag", limit.Tag)
	}
}

func TestConcurrencyLimitsService_Get(t *testing.T) {
	id := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %v, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ConcurrencyLimit{ID: id, Tag: "test"})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	limit, err := client.ConcurrencyLimits.Get(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limit.ID != id {
		t.Errorf("ID = %v, want %v", limit.ID, id)
	}
}

func TestConcurrencyLimitsService_GetByTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %v, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ConcurrencyLimit{ID: uuid.New(), Tag: "my-tag"})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	limit, err := client.ConcurrencyLimits.GetByTag(ctx, "my-tag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limit.Tag != "my-tag" {
		t.Errorf("Tag = %v, want my-tag", limit.Tag)
	}
}

func TestConcurrencyLimitsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	err := client.ConcurrencyLimits.Delete(ctx, uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConcurrencyLimitsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.ConcurrencyLimit{{ID: uuid.New(), Tag: "test"}})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	limits, err := client.ConcurrencyLimits.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(limits) != 1 {
		t.Errorf("count = %v, want 1", len(limits))
	}
}
