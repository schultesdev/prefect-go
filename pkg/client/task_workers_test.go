package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
)

func TestTaskWorkersService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/task_workers/filter" {
			t.Errorf("path = %v, want /api/task_workers/filter", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.TaskWorkerResponse{
			{Identifier: "worker-1", TaskKeys: []string{"key1"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	workers, err := client.TaskWorkers.List(ctx, &models.TaskWorkerFilter{TaskKeys: []string{"key1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(workers) != 1 {
		t.Errorf("count = %v, want 1", len(workers))
	}
	if workers[0].Identifier != "worker-1" {
		t.Errorf("Identifier = %v, want worker-1", workers[0].Identifier)
	}
}
