package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// WorkQueuesService handles operations related to work queues.
type WorkQueuesService struct {
	client *Client
}

// Get retrieves a work queue by ID.
func (w *WorkQueuesService) Get(ctx context.Context, id uuid.UUID) (*models.WorkQueue, error) {
	var workQueue models.WorkQueue
	path := joinPath("/work_queues", id.String())
	if err := w.client.get(ctx, path, &workQueue); err != nil {
		return nil, fmt.Errorf("failed to get work queue: %w", err)
	}
	return &workQueue, nil
}

// Create creates a new work queue.
func (w *WorkQueuesService) Create(ctx context.Context, req *models.WorkQueueCreate) (*models.WorkQueue, error) {
	var workQueue models.WorkQueue
	if err := w.client.post(ctx, "/work_queues/", req, &workQueue); err != nil {
		return nil, fmt.Errorf("failed to create work queue: %w", err)
	}
	return &workQueue, nil
}

// GetByName retrieves a work queue by name.
func (w *WorkQueuesService) GetByName(ctx context.Context, name string) (*models.WorkQueue, error) {
	var workQueue models.WorkQueue
	path := joinPath("/work_queues/name", name)
	if err := w.client.get(ctx, path, &workQueue); err != nil {
		return nil, fmt.Errorf("failed to get work queue by name: %w", err)
	}
	return &workQueue, nil
}

// Update updates a work queue.
func (w *WorkQueuesService) Update(ctx context.Context, id uuid.UUID, req *models.WorkQueueUpdate) error {
	path := joinPath("/work_queues", id.String())
	if err := w.client.patch(ctx, path, req, nil); err != nil {
		return fmt.Errorf("failed to update work queue: %w", err)
	}
	return nil
}

// Delete deletes a work queue by ID.
func (w *WorkQueuesService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/work_queues", id.String())
	if err := w.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete work queue: %w", err)
	}
	return nil
}

// List retrieves work queues with optional filtering.
func (w *WorkQueuesService) List(ctx context.Context, filter *models.WorkQueueFilter, offset, limit int) ([]models.WorkQueue, error) {
	if filter == nil {
		filter = &models.WorkQueueFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	var workQueues []models.WorkQueue
	if err := w.client.post(ctx, "/work_queues/filter", filter, &workQueues); err != nil {
		return nil, fmt.Errorf("failed to list work queues: %w", err)
	}
	return workQueues, nil
}
