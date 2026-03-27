package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
)

// TaskWorkersService handles operations related to task workers.
type TaskWorkersService struct {
	client *Client
}

// List retrieves task workers matching the filter.
func (s *TaskWorkersService) List(ctx context.Context, filter *models.TaskWorkerFilter) ([]models.TaskWorkerResponse, error) {
	if filter == nil {
		filter = &models.TaskWorkerFilter{}
	}

	var workers []models.TaskWorkerResponse
	if err := s.client.post(ctx, "/task_workers/filter", filter, &workers); err != nil {
		return nil, fmt.Errorf("failed to list task workers: %w", err)
	}
	return workers, nil
}
