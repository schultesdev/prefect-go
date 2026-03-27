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

func TestArtifactsService_Create(t *testing.T) {
	expected := models.Artifact{
		ID:  uuid.New(),
		Key: strPtr("test-key"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/artifacts/" {
			t.Errorf("path = %v, want /api/artifacts/", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	artifact, err := client.Artifacts.Create(ctx, &models.ArtifactCreate{Key: strPtr("test-key")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if artifact.ID != expected.ID {
		t.Errorf("ID = %v, want %v", artifact.ID, expected.ID)
	}
}

func TestArtifactsService_Get(t *testing.T) {
	id := uuid.New()
	expected := models.Artifact{ID: id, Key: strPtr("test-key")}

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

	artifact, err := client.Artifacts.Get(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if artifact.ID != id {
		t.Errorf("ID = %v, want %v", artifact.ID, id)
	}
}

func TestArtifactsService_Delete(t *testing.T) {
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

	err := client.Artifacts.Delete(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestArtifactsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		if r.URL.Path != "/api/artifacts/filter" {
			t.Errorf("path = %v, want /api/artifacts/filter", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]interface{}{
			"results": []models.Artifact{{ID: uuid.New()}},
			"count":   1,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	result, err := client.Artifacts.List(ctx, nil, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Results) != 1 {
		t.Errorf("results count = %v, want 1", len(result.Results))
	}
}

func TestArtifactsService_Count(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`5`))
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	count, err := client.Artifacts.Count(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 5 {
		t.Errorf("count = %v, want 5", count)
	}
}

func strPtr(s string) *string {
	return &s
}
