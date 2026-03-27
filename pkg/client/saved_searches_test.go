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

func TestSavedSearchesService_Create(t *testing.T) {
	expected := models.SavedSearch{ID: uuid.New(), Name: "test-search"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %v, want PUT", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	search, err := client.SavedSearches.Create(ctx, &models.SavedSearchCreate{Name: "test-search"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if search.Name != expected.Name {
		t.Errorf("Name = %v, want %v", search.Name, expected.Name)
	}
}

func TestSavedSearchesService_Get(t *testing.T) {
	id := uuid.New()
	expected := models.SavedSearch{ID: id, Name: "test"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %v, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	search, err := client.SavedSearches.Get(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if search.ID != id {
		t.Errorf("ID = %v, want %v", search.ID, id)
	}
}

func TestSavedSearchesService_Delete(t *testing.T) {
	id := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	err := client.SavedSearches.Delete(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSavedSearchesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/saved_searches/filter" {
			t.Errorf("path = %v, want /api/saved_searches/filter", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.SavedSearch{{ID: uuid.New(), Name: "test"}})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	searches, err := client.SavedSearches.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(searches) != 1 {
		t.Errorf("count = %v, want 1", len(searches))
	}
}
