package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// FlowRunStatesService handles operations related to flow run states.
type FlowRunStatesService struct {
	client *Client
}

// List retrieves all flow run states.
func (s *FlowRunStatesService) List(ctx context.Context) ([]models.State, error) {
	var states []models.State
	if err := s.client.get(ctx, "/flow_run_states/", &states); err != nil {
		return nil, fmt.Errorf("failed to list flow run states: %w", err)
	}
	return states, nil
}

// Get retrieves a flow run state by ID.
func (s *FlowRunStatesService) Get(ctx context.Context, id uuid.UUID) (*models.State, error) {
	var state models.State
	path := joinPath("/flow_run_states", id.String())
	if err := s.client.get(ctx, path, &state); err != nil {
		return nil, fmt.Errorf("failed to get flow run state: %w", err)
	}
	return &state, nil
}
