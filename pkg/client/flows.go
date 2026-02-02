package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gregor/prefect-go/pkg/models"
	"github.com/gregor/prefect-go/pkg/pagination"
)

// FlowsService handles operations related to flows.
type FlowsService struct {
	client *Client
}

// Create creates a new flow.
func (s *FlowsService) Create(ctx context.Context, req *models.FlowCreate) (*models.Flow, error) {
	var flow models.Flow
	if err := s.client.post(ctx, "/flows/", req, &flow); err != nil {
		return nil, fmt.Errorf("failed to create flow: %w", err)
	}
	return &flow, nil
}

// Get retrieves a flow by ID.
func (s *FlowsService) Get(ctx context.Context, id uuid.UUID) (*models.Flow, error) {
	var flow models.Flow
	path := joinPath("/flows", id.String())
	if err := s.client.get(ctx, path, &flow); err != nil {
		return nil, fmt.Errorf("failed to get flow: %w", err)
	}
	return &flow, nil
}

// GetByName retrieves a flow by name.
func (s *FlowsService) GetByName(ctx context.Context, name string) (*models.Flow, error) {
	var flow models.Flow
	path := buildPath("/flows/name/"+name, nil)
	if err := s.client.get(ctx, path, &flow); err != nil {
		return nil, fmt.Errorf("failed to get flow by name: %w", err)
	}
	return &flow, nil
}

// List retrieves a list of flows with optional filtering.
func (s *FlowsService) List(ctx context.Context, filter *models.FlowFilter, offset, limit int) (*pagination.PaginatedResponse[models.Flow], error) {
	if filter == nil {
		filter = &models.FlowFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	type response struct {
		Results []models.Flow `json:"results"`
		Count   int           `json:"count"`
	}

	var resp response
	if err := s.client.post(ctx, "/flows/filter", filter, &resp); err != nil {
		return nil, fmt.Errorf("failed to list flows: %w", err)
	}

	return &pagination.PaginatedResponse[models.Flow]{
		Results: resp.Results,
		Count:   resp.Count,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(resp.Results) < resp.Count,
	}, nil
}

// ListAll returns an iterator for all flows matching the filter.
func (s *FlowsService) ListAll(ctx context.Context, filter *models.FlowFilter) *pagination.Iterator[models.Flow] {
	fetchFunc := func(ctx context.Context, offset, limit int) (*pagination.PaginatedResponse[models.Flow], error) {
		return s.List(ctx, filter, offset, limit)
	}
	return pagination.NewIterator(fetchFunc, 100)
}

// Update updates a flow.
func (s *FlowsService) Update(ctx context.Context, id uuid.UUID, req *models.FlowUpdate) (*models.Flow, error) {
	var flow models.Flow
	path := joinPath("/flows", id.String())
	if err := s.client.patch(ctx, path, req, &flow); err != nil {
		return nil, fmt.Errorf("failed to update flow: %w", err)
	}
	return &flow, nil
}

// Delete deletes a flow by ID.
func (s *FlowsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/flows", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete flow: %w", err)
	}
	return nil
}

// Count returns the number of flows matching the filter.
func (s *FlowsService) Count(ctx context.Context, filter *models.FlowFilter) (int, error) {
	if filter == nil {
		filter = &models.FlowFilter{}
	}

	var result struct {
		Count int `json:"count"`
	}

	if err := s.client.post(ctx, "/flows/count", filter, &result); err != nil {
		return 0, fmt.Errorf("failed to count flows: %w", err)
	}

	return result.Count, nil
}
