package client

import (
	"context"
	"fmt"
	"net/http"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// SavedSearchesService handles operations related to saved searches.
type SavedSearchesService struct {
	client *Client
}

// Create creates a new saved search.
func (s *SavedSearchesService) Create(ctx context.Context, req *models.SavedSearchCreate) (*models.SavedSearch, error) {
	var savedSearch models.SavedSearch
	if err := s.client.do(ctx, http.MethodPut, "/saved_searches/", req, &savedSearch); err != nil {
		return nil, fmt.Errorf("failed to create saved search: %w", err)
	}
	return &savedSearch, nil
}

// Get retrieves a saved search by ID.
func (s *SavedSearchesService) Get(ctx context.Context, id uuid.UUID) (*models.SavedSearch, error) {
	var savedSearch models.SavedSearch
	path := joinPath("/saved_searches", id.String())
	if err := s.client.get(ctx, path, &savedSearch); err != nil {
		return nil, fmt.Errorf("failed to get saved search: %w", err)
	}
	return &savedSearch, nil
}

// Delete deletes a saved search by ID.
func (s *SavedSearchesService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/saved_searches", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete saved search: %w", err)
	}
	return nil
}

// List retrieves all saved searches.
func (s *SavedSearchesService) List(ctx context.Context) ([]models.SavedSearch, error) {
	var savedSearches []models.SavedSearch
	if err := s.client.post(ctx, "/saved_searches/filter", nil, &savedSearches); err != nil {
		return nil, fmt.Errorf("failed to list saved searches: %w", err)
	}
	return savedSearches, nil
}
