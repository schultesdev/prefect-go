package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// TaskRunStatesService handles operations related to task run states.
type TaskRunStatesService struct {
	client *Client
}

// List retrieves all task run states.
func (s *TaskRunStatesService) List(ctx context.Context) ([]models.State, error) {
	var states []models.State
	if err := s.client.get(ctx, "/task_run_states/", &states); err != nil {
		return nil, fmt.Errorf("failed to list task run states: %w", err)
	}
	return states, nil
}

// Get retrieves a task run state by ID.
func (s *TaskRunStatesService) Get(ctx context.Context, id uuid.UUID) (*models.State, error) {
	var state models.State
	path := joinPath("/task_run_states", id.String())
	if err := s.client.get(ctx, path, &state); err != nil {
		return nil, fmt.Errorf("failed to get task run state: %w", err)
	}
	return &state, nil
}
