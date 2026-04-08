package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"git.schultes.dev/schultesdev/prefect-go/pkg/pagination"
	"github.com/google/uuid"
)

// TaskRunsService handles operations related to task runs.
type TaskRunsService struct {
	client *Client
}

// Create creates a new task run.
func (t *TaskRunsService) Create(ctx context.Context, req *models.TaskRunCreate) (*models.TaskRun, error) {
	var taskRun models.TaskRun
	if err := t.client.post(ctx, "/task_runs/", req, &taskRun); err != nil {
		return nil, fmt.Errorf("failed to create task run: %w", err)
	}
	return &taskRun, nil
}

// Get retrieves a task run by ID.
func (t *TaskRunsService) Get(ctx context.Context, id uuid.UUID) (*models.TaskRun, error) {
	var taskRun models.TaskRun
	path := joinPath("/task_runs", id.String())
	if err := t.client.get(ctx, path, &taskRun); err != nil {
		return nil, fmt.Errorf("failed to get task run: %w", err)
	}
	return &taskRun, nil
}

// Delete deletes a task run by ID.
func (t *TaskRunsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/task_runs", id.String())
	if err := t.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete task run: %w", err)
	}
	return nil
}

// SetState sets the state of a task run.
func (t *TaskRunsService) SetState(ctx context.Context, id uuid.UUID, state *models.StateCreate) (*models.TaskRun, error) {
	var taskRun models.TaskRun
	path := joinPath("/task_runs", id.String(), "set_state")

	req := struct {
		State *models.StateCreate `json:"state"`
	}{
		State: state,
	}

	if err := t.client.post(ctx, path, req, &taskRun); err != nil {
		return nil, fmt.Errorf("failed to set task run state: %w", err)
	}
	return &taskRun, nil
}

// Update updates a task run.
func (t *TaskRunsService) Update(ctx context.Context, id uuid.UUID, req *models.TaskRunUpdate) (*models.TaskRun, error) {
	var taskRun models.TaskRun
	path := joinPath("/task_runs", id.String())
	if err := t.client.patch(ctx, path, req, &taskRun); err != nil {
		return nil, fmt.Errorf("failed to update task run: %w", err)
	}
	return &taskRun, nil
}

// List retrieves a list of task runs with optional filtering.
func (t *TaskRunsService) List(ctx context.Context, filter *models.TaskRunFilter, offset, limit int) (*pagination.PaginatedResponse[models.TaskRun], error) {
	if filter == nil {
		filter = &models.TaskRunFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	type response struct {
		Results []models.TaskRun `json:"results"`
		Count   int              `json:"count"`
	}

	var resp response
	if err := t.client.post(ctx, "/task_runs/filter", filter, &resp); err != nil {
		return nil, fmt.Errorf("failed to list task runs: %w", err)
	}

	return &pagination.PaginatedResponse[models.TaskRun]{
		Results: resp.Results,
		Count:   resp.Count,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(resp.Results) < resp.Count,
	}, nil
}

// ListAll returns an iterator for all task runs matching the filter.
func (t *TaskRunsService) ListAll(ctx context.Context, filter *models.TaskRunFilter) *pagination.Iterator[models.TaskRun] {
	fetchFunc := func(ctx context.Context, offset, limit int) (*pagination.PaginatedResponse[models.TaskRun], error) {
		return t.List(ctx, filter, offset, limit)
	}
	return pagination.NewIterator(fetchFunc, 100)
}

// Count returns the number of task runs matching the filter.
func (t *TaskRunsService) Count(ctx context.Context, filter *models.TaskRunFilter) (int, error) {
	if filter == nil {
		filter = &models.TaskRunFilter{}
	}

	var count int
	if err := t.client.post(ctx, "/task_runs/count", filter, &count); err != nil {
		return 0, fmt.Errorf("failed to count task runs: %w", err)
	}
	return count, nil
}
