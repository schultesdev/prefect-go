package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"git.schultes.dev/schultesdev/prefect-go/pkg/pagination"
	"github.com/google/uuid"
)

// ArtifactsService handles operations related to artifacts.
type ArtifactsService struct {
	client *Client
}

// Create creates a new artifact.
func (s *ArtifactsService) Create(ctx context.Context, req *models.ArtifactCreate) (*models.Artifact, error) {
	var artifact models.Artifact
	if err := s.client.post(ctx, "/artifacts/", req, &artifact); err != nil {
		return nil, fmt.Errorf("failed to create artifact: %w", err)
	}
	return &artifact, nil
}

// Get retrieves an artifact by ID.
func (s *ArtifactsService) Get(ctx context.Context, id uuid.UUID) (*models.Artifact, error) {
	var artifact models.Artifact
	path := joinPath("/artifacts", id.String())
	if err := s.client.get(ctx, path, &artifact); err != nil {
		return nil, fmt.Errorf("failed to get artifact: %w", err)
	}
	return &artifact, nil
}

// GetLatest retrieves the latest artifact for a given key.
func (s *ArtifactsService) GetLatest(ctx context.Context, key string) (*models.Artifact, error) {
	var artifact models.Artifact
	path := joinPath("/artifacts", key, "latest")
	if err := s.client.get(ctx, path, &artifact); err != nil {
		return nil, fmt.Errorf("failed to get latest artifact: %w", err)
	}
	return &artifact, nil
}

// Update updates an artifact.
func (s *ArtifactsService) Update(ctx context.Context, id uuid.UUID, req *models.ArtifactUpdate) error {
	path := joinPath("/artifacts", id.String())
	if err := s.client.patch(ctx, path, req, nil); err != nil {
		return fmt.Errorf("failed to update artifact: %w", err)
	}
	return nil
}

// Delete deletes an artifact by ID.
func (s *ArtifactsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/artifacts", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete artifact: %w", err)
	}
	return nil
}

// List retrieves a list of artifacts with optional filtering.
func (s *ArtifactsService) List(ctx context.Context, filter *models.ArtifactFilter, offset, limit int) (*pagination.PaginatedResponse[models.Artifact], error) {
	if filter == nil {
		filter = &models.ArtifactFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	type response struct {
		Results []models.Artifact `json:"results"`
		Count   int               `json:"count"`
	}

	var resp response
	if err := s.client.post(ctx, "/artifacts/filter", filter, &resp); err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}

	return &pagination.PaginatedResponse[models.Artifact]{
		Results: resp.Results,
		Count:   resp.Count,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(resp.Results) < resp.Count,
	}, nil
}

// ListAll returns an iterator for all artifacts matching the filter.
func (s *ArtifactsService) ListAll(ctx context.Context, filter *models.ArtifactFilter) *pagination.Iterator[models.Artifact] {
	fetchFunc := func(ctx context.Context, offset, limit int) (*pagination.PaginatedResponse[models.Artifact], error) {
		return s.List(ctx, filter, offset, limit)
	}
	return pagination.NewIterator(fetchFunc, 100)
}

// Count returns the number of artifacts matching the filter.
func (s *ArtifactsService) Count(ctx context.Context, filter *models.ArtifactFilter) (int, error) {
	if filter == nil {
		filter = &models.ArtifactFilter{}
	}

	var count int
	if err := s.client.post(ctx, "/artifacts/count", filter, &count); err != nil {
		return 0, fmt.Errorf("failed to count artifacts: %w", err)
	}
	return count, nil
}

// ListLatest retrieves the latest artifact collections with optional filtering.
func (s *ArtifactsService) ListLatest(ctx context.Context, filter *models.ArtifactCollectionFilter, offset, limit int) (*pagination.PaginatedResponse[models.ArtifactCollection], error) {
	if filter == nil {
		filter = &models.ArtifactCollectionFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	type response struct {
		Results []models.ArtifactCollection `json:"results"`
		Count   int                         `json:"count"`
	}

	var resp response
	if err := s.client.post(ctx, "/artifacts/latest/filter", filter, &resp); err != nil {
		return nil, fmt.Errorf("failed to list latest artifacts: %w", err)
	}

	return &pagination.PaginatedResponse[models.ArtifactCollection]{
		Results: resp.Results,
		Count:   resp.Count,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(resp.Results) < resp.Count,
	}, nil
}

// CountLatest returns the number of latest artifact collections matching the filter.
func (s *ArtifactsService) CountLatest(ctx context.Context, filter *models.ArtifactCollectionFilter) (int, error) {
	if filter == nil {
		filter = &models.ArtifactCollectionFilter{}
	}

	var count int
	if err := s.client.post(ctx, "/artifacts/latest/count", filter, &count); err != nil {
		return 0, fmt.Errorf("failed to count latest artifacts: %w", err)
	}
	return count, nil
}
