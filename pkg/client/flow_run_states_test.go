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

func TestFlowRunStatesService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/flow_run_states/" {
			t.Errorf("path = %v, want /api/flow_run_states/", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %v, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.State{
			{ID: uuid.New(), Type: models.StateTypeCompleted},
		})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	states, err := client.FlowRunStates.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(states) != 1 {
		t.Errorf("count = %v, want 1", len(states))
	}
}

func TestFlowRunStatesService_Get(t *testing.T) {
	id := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %v, want GET", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.State{ID: id, Type: models.StateTypeRunning})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	state, err := client.FlowRunStates.Get(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.ID != id {
		t.Errorf("ID = %v, want %v", state.ID, id)
	}
}
