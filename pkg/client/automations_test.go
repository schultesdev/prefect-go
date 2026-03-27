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

func TestAutomationsService_Create(t *testing.T) {
	expected := models.Automation{
		ID:   uuid.New(),
		Name: "test-automation",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/automations/" {
			t.Errorf("path = %v, want /api/automations/", r.URL.Path)
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

	automation, err := client.Automations.Create(ctx, &models.AutomationCreate{
		Name:    "test-automation",
		Trigger: json.RawMessage(`{"type":"event"}`),
		Actions: []json.RawMessage{json.RawMessage(`{"type":"do-nothing"}`)},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if automation.Name != expected.Name {
		t.Errorf("Name = %v, want %v", automation.Name, expected.Name)
	}
}

func TestAutomationsService_Get(t *testing.T) {
	id := uuid.New()
	expected := models.Automation{ID: id, Name: "test"}

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

	automation, err := client.Automations.Get(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if automation.ID != id {
		t.Errorf("ID = %v, want %v", automation.ID, id)
	}
}

func TestAutomationsService_Delete(t *testing.T) {
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

	err := client.Automations.Delete(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAutomationsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/automations/filter" {
			t.Errorf("path = %v, want /api/automations/filter", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Automation{{ID: uuid.New(), Name: "test"}})
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	automations, err := client.Automations.List(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(automations) != 1 {
		t.Errorf("count = %v, want 1", len(automations))
	}
}

func TestAutomationsService_Patch(t *testing.T) {
	id := uuid.New()
	enabled := true

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %v, want PATCH", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))
	ctx := context.Background()

	err := client.Automations.Patch(ctx, id, &models.AutomationPartialUpdate{Enabled: &enabled})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
