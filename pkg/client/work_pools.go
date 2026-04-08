package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
)

// WorkPoolsService handles operations related to work pools.
type WorkPoolsService struct {
	client *Client
}

// Get retrieves a work pool by name.
func (w *WorkPoolsService) Get(ctx context.Context, name string) (*models.WorkPool, error) {
	var workPool models.WorkPool
	path := joinPath("/work_pools", name)
	if err := w.client.get(ctx, path, &workPool); err != nil {
		return nil, fmt.Errorf("failed to get work pool: %w", err)
	}
	return &workPool, nil
}

// Create creates a new work pool.
func (w *WorkPoolsService) Create(ctx context.Context, req *models.WorkPoolCreate) (*models.WorkPool, error) {
	var workPool models.WorkPool
	if err := w.client.post(ctx, "/work_pools/", req, &workPool); err != nil {
		return nil, fmt.Errorf("failed to create work pool: %w", err)
	}
	return &workPool, nil
}

// Update updates a work pool.
func (w *WorkPoolsService) Update(ctx context.Context, name string, req *models.WorkPoolUpdate) error {
	path := joinPath("/work_pools", name)
	if err := w.client.patch(ctx, path, req, nil); err != nil {
		return fmt.Errorf("failed to update work pool: %w", err)
	}
	return nil
}

// Delete deletes a work pool by name.
func (w *WorkPoolsService) Delete(ctx context.Context, name string) error {
	path := joinPath("/work_pools", name)
	if err := w.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete work pool: %w", err)
	}
	return nil
}

// List retrieves work pools with optional filtering.
func (w *WorkPoolsService) List(ctx context.Context, filter *models.WorkPoolFilter, offset, limit int) ([]models.WorkPool, error) {
	if filter == nil {
		filter = &models.WorkPoolFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	var workPools []models.WorkPool
	if err := w.client.post(ctx, "/work_pools/filter", filter, &workPools); err != nil {
		return nil, fmt.Errorf("failed to list work pools: %w", err)
	}
	return workPools, nil
}

// Count returns the number of work pools matching the filter.
func (w *WorkPoolsService) Count(ctx context.Context, filter *models.WorkPoolFilter) (int, error) {
	if filter == nil {
		filter = &models.WorkPoolFilter{}
	}

	var count int
	if err := w.client.post(ctx, "/work_pools/count", filter, &count); err != nil {
		return 0, fmt.Errorf("failed to count work pools: %w", err)
	}
	return count, nil
}
