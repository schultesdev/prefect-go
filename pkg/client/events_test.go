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

func TestEventsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/events/filter" {
			t.Errorf("path = %v, want /api/events/filter", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.EventPage{
			Events: []models.Event{{ID: uuid.New(), Event: "test.event"}},
			Total:  1,
		})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	page, err := client.Events.List(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Total != 1 {
		t.Errorf("total = %v, want 1", page.Total)
	}
}

func TestEventsService_CountBy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.EventCount{{Value: "test", Count: 5}})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	counts, err := client.Events.CountBy(ctx, "event", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(counts) != 1 {
		t.Errorf("count = %v, want 1", len(counts))
	}
}
