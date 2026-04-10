package client

import (
	"context"
	"fmt"
	"time"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"git.schultes.dev/schultesdev/prefect-go/pkg/pagination"
	"github.com/google/uuid"
)

// FlowRunsService handles operations related to flow runs.
type FlowRunsService struct {
	client *Client
}

// Create creates a new flow run.
func (s *FlowRunsService) Create(ctx context.Context, req *models.FlowRunCreate) (*models.FlowRun, error) {
	var flowRun models.FlowRun
	if err := s.client.post(ctx, "/flow_runs/", req, &flowRun); err != nil {
		return nil, fmt.Errorf("failed to create flow run: %w", err)
	}
	return &flowRun, nil
}

// Get retrieves a flow run by ID.
func (s *FlowRunsService) Get(ctx context.Context, id uuid.UUID) (*models.FlowRun, error) {
	var flowRun models.FlowRun
	path := joinPath("/flow_runs", id.String())
	if err := s.client.get(ctx, path, &flowRun); err != nil {
		return nil, fmt.Errorf("failed to get flow run: %w", err)
	}
	return &flowRun, nil
}

// List retrieves a list of flow runs with optional filtering.
func (s *FlowRunsService) List(ctx context.Context, filter *models.FlowRunFilter, offset, limit int) (*pagination.PaginatedResponse[models.FlowRun], error) {
	if filter == nil {
		filter = &models.FlowRunFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	var results []models.FlowRun
	if err := s.client.post(ctx, "/flow_runs/filter", filter, &results); err != nil {
		return nil, fmt.Errorf("failed to list flow runs: %w", err)
	}

	return &pagination.PaginatedResponse[models.FlowRun]{
		Results: results,
		Count:   offset + len(results),
		Limit:   limit,
		Offset:  offset,
		HasMore: len(results) == limit,
	}, nil
}

// ListAll returns an iterator for all flow runs matching the filter.
func (s *FlowRunsService) ListAll(ctx context.Context, filter *models.FlowRunFilter) *pagination.Iterator[models.FlowRun] {
	fetchFunc := func(ctx context.Context, offset, limit int) (*pagination.PaginatedResponse[models.FlowRun], error) {
		return s.List(ctx, filter, offset, limit)
	}
	return pagination.NewIterator(fetchFunc, 100)
}

// Update updates a flow run.
func (s *FlowRunsService) Update(ctx context.Context, id uuid.UUID, req *models.FlowRunUpdate) (*models.FlowRun, error) {
	var flowRun models.FlowRun
	path := joinPath("/flow_runs", id.String())
	if err := s.client.patch(ctx, path, req, &flowRun); err != nil {
		return nil, fmt.Errorf("failed to update flow run: %w", err)
	}
	return &flowRun, nil
}

// Delete deletes a flow run by ID.
func (s *FlowRunsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/flow_runs", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete flow run: %w", err)
	}
	return nil
}

// SetState sets the state of a flow run.
func (s *FlowRunsService) SetState(ctx context.Context, id uuid.UUID, state *models.StateCreate) (*models.FlowRun, error) {
	var flowRun models.FlowRun
	path := joinPath("/flow_runs", id.String(), "set_state")

	req := struct {
		State *models.StateCreate `json:"state"`
	}{
		State: state,
	}

	if err := s.client.post(ctx, path, req, &flowRun); err != nil {
		return nil, fmt.Errorf("failed to set flow run state: %w", err)
	}
	return &flowRun, nil
}

// Wait waits for a flow run to complete by polling its state.
// It returns the final flow run state or an error if the context is cancelled.
func (s *FlowRunsService) Wait(ctx context.Context, id uuid.UUID, pollInterval time.Duration) (*models.FlowRun, error) {
	if pollInterval <= 0 {
		pollInterval = 5 * time.Second
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			flowRun, err := s.Get(ctx, id)
			if err != nil {
				return nil, err
			}

			// Check if flow run is in a terminal state
			if flowRun.StateType != nil {
				switch *flowRun.StateType {
				case models.StateTypeCompleted,
					models.StateTypeFailed,
					models.StateTypeCancelled,
					models.StateTypeCrashed:
					return flowRun, nil
				}
			}
		}
	}
}

// Count returns the number of flow runs matching the filter.
func (s *FlowRunsService) Count(ctx context.Context, filter *models.FlowRunFilter) (int, error) {
	if filter == nil {
		filter = &models.FlowRunFilter{}
	}

	var result struct {
		Count int `json:"count"`
	}

	if err := s.client.post(ctx, "/flow_runs/count", filter, &result); err != nil {
		return 0, fmt.Errorf("failed to count flow runs: %w", err)
	}

	return result.Count, nil
}
