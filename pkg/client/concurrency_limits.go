package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// ConcurrencyLimitsService handles operations related to v1 concurrency limits.
type ConcurrencyLimitsService struct {
	client *Client
}

// Create creates a new concurrency limit.
func (s *ConcurrencyLimitsService) Create(ctx context.Context, req *models.ConcurrencyLimitCreate) (*models.ConcurrencyLimit, error) {
	var limit models.ConcurrencyLimit
	if err := s.client.post(ctx, "/concurrency_limits/", req, &limit); err != nil {
		return nil, fmt.Errorf("failed to create concurrency limit: %w", err)
	}
	return &limit, nil
}

// Get retrieves a concurrency limit by ID.
func (s *ConcurrencyLimitsService) Get(ctx context.Context, id uuid.UUID) (*models.ConcurrencyLimit, error) {
	var limit models.ConcurrencyLimit
	path := joinPath("/concurrency_limits", id.String())
	if err := s.client.get(ctx, path, &limit); err != nil {
		return nil, fmt.Errorf("failed to get concurrency limit: %w", err)
	}
	return &limit, nil
}

// GetByTag retrieves a concurrency limit by tag.
func (s *ConcurrencyLimitsService) GetByTag(ctx context.Context, tag string) (*models.ConcurrencyLimit, error) {
	var limit models.ConcurrencyLimit
	path := joinPath("/concurrency_limits/tag", tag)
	if err := s.client.get(ctx, path, &limit); err != nil {
		return nil, fmt.Errorf("failed to get concurrency limit by tag: %w", err)
	}
	return &limit, nil
}

// Delete deletes a concurrency limit by ID.
func (s *ConcurrencyLimitsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/concurrency_limits", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete concurrency limit: %w", err)
	}
	return nil
}

// DeleteByTag deletes a concurrency limit by tag.
func (s *ConcurrencyLimitsService) DeleteByTag(ctx context.Context, tag string) error {
	path := joinPath("/concurrency_limits/tag", tag)
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete concurrency limit by tag: %w", err)
	}
	return nil
}

// ResetByTag resets a concurrency limit by tag.
func (s *ConcurrencyLimitsService) ResetByTag(ctx context.Context, tag string) error {
	path := joinPath("/concurrency_limits/tag", tag, "reset")
	if err := s.client.post(ctx, path, nil, nil); err != nil {
		return fmt.Errorf("failed to reset concurrency limit: %w", err)
	}
	return nil
}

// List retrieves concurrency limits with optional filtering.
func (s *ConcurrencyLimitsService) List(ctx context.Context) ([]models.ConcurrencyLimit, error) {
	var limits []models.ConcurrencyLimit
	if err := s.client.post(ctx, "/concurrency_limits/filter", nil, &limits); err != nil {
		return nil, fmt.Errorf("failed to list concurrency limits: %w", err)
	}
	return limits, nil
}

// Increment increments concurrency limit slots.
func (s *ConcurrencyLimitsService) Increment(ctx context.Context, names []string, slots int, mode string) ([]models.MinimalConcurrencyLimitResponse, error) {
	req := struct {
		Names []string `json:"names"`
		Slots int      `json:"slots"`
		Mode  string   `json:"mode,omitempty"`
	}{
		Names: names,
		Slots: slots,
		Mode:  mode,
	}

	var resp []models.MinimalConcurrencyLimitResponse
	if err := s.client.post(ctx, "/concurrency_limits/increment", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to increment concurrency limits: %w", err)
	}
	return resp, nil
}

// Decrement decrements concurrency limit slots.
func (s *ConcurrencyLimitsService) Decrement(ctx context.Context, names []string, slots int) error {
	req := struct {
		Names []string `json:"names"`
		Slots int      `json:"slots"`
	}{
		Names: names,
		Slots: slots,
	}

	if err := s.client.post(ctx, "/concurrency_limits/decrement", req, nil); err != nil {
		return fmt.Errorf("failed to decrement concurrency limits: %w", err)
	}
	return nil
}
