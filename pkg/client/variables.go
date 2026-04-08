package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"git.schultes.dev/schultesdev/prefect-go/pkg/pagination"
	"github.com/google/uuid"
)

// VariablesService handles operations related to variables.
type VariablesService struct {
	client *Client
}

// Create creates a new variable.
func (v *VariablesService) Create(ctx context.Context, req *models.VariableCreate) (*models.Variable, error) {
	var variable models.Variable
	if err := v.client.post(ctx, "/variables/", req, &variable); err != nil {
		return nil, fmt.Errorf("failed to create variable: %w", err)
	}
	return &variable, nil
}

// Get retrieves a variable by ID.
func (v *VariablesService) Get(ctx context.Context, id uuid.UUID) (*models.Variable, error) {
	var variable models.Variable
	path := joinPath("/variables", id.String())
	if err := v.client.get(ctx, path, &variable); err != nil {
		return nil, fmt.Errorf("failed to get variable: %w", err)
	}
	return &variable, nil
}

// GetByName retrieves a variable by name.
func (v *VariablesService) GetByName(ctx context.Context, name string) (*models.Variable, error) {
	var variable models.Variable
	path := joinPath("/variables/name", name)
	if err := v.client.get(ctx, path, &variable); err != nil {
		return nil, fmt.Errorf("failed to get variable by name: %w", err)
	}
	return &variable, nil
}

// Update updates a variable.
func (v *VariablesService) Update(ctx context.Context, id uuid.UUID, req *models.VariableUpdate) (*models.Variable, error) {
	var variable models.Variable
	path := joinPath("/variables", id.String())
	if err := v.client.patch(ctx, path, req, &variable); err != nil {
		return nil, fmt.Errorf("failed to update variable: %w", err)
	}
	return &variable, nil
}

// Delete deletes a variable by ID.
func (v *VariablesService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/variables", id.String())
	if err := v.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete variable: %w", err)
	}
	return nil
}

// List retrieves a list of variables with optional filtering.
func (v *VariablesService) List(ctx context.Context, filter *models.VariableFilter, offset, limit int) (*pagination.PaginatedResponse[models.Variable], error) {
	if filter == nil {
		filter = &models.VariableFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	type response struct {
		Results []models.Variable `json:"results"`
		Count   int               `json:"count"`
	}

	var resp response
	if err := v.client.post(ctx, "/variables/filter", filter, &resp); err != nil {
		return nil, fmt.Errorf("failed to list variables: %w", err)
	}

	return &pagination.PaginatedResponse[models.Variable]{
		Results: resp.Results,
		Count:   resp.Count,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(resp.Results) < resp.Count,
	}, nil
}

// ListAll returns an iterator for all variables matching the filter.
func (v *VariablesService) ListAll(ctx context.Context, filter *models.VariableFilter) *pagination.Iterator[models.Variable] {
	fetchFunc := func(ctx context.Context, offset, limit int) (*pagination.PaginatedResponse[models.Variable], error) {
		return v.List(ctx, filter, offset, limit)
	}
	return pagination.NewIterator(fetchFunc, 100)
}

// Count returns the number of variables matching the filter.
func (v *VariablesService) Count(ctx context.Context, filter *models.VariableFilter) (int, error) {
	if filter == nil {
		filter = &models.VariableFilter{}
	}

	var count int
	if err := v.client.post(ctx, "/variables/count", filter, &count); err != nil {
		return 0, fmt.Errorf("failed to count variables: %w", err)
	}
	return count, nil
}
