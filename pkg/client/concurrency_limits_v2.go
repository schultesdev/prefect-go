package client

import (
	"context"
	"fmt"
	"net/http"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// ConcurrencyLimitsV2Service handles operations related to v2 concurrency limits.
type ConcurrencyLimitsV2Service struct {
	client *Client
}

// Create creates a new v2 concurrency limit.
func (s *ConcurrencyLimitsV2Service) Create(ctx context.Context, req *models.ConcurrencyLimitV2Create) (*models.ConcurrencyLimitV2, error) {
	var limit models.ConcurrencyLimitV2
	if err := s.client.do(ctx, http.MethodPost, "/v2/concurrency_limits/", req, &limit); err != nil {
		return nil, fmt.Errorf("failed to create v2 concurrency limit: %w", err)
	}
	return &limit, nil
}

// Get retrieves a v2 concurrency limit by ID or name.
func (s *ConcurrencyLimitsV2Service) Get(ctx context.Context, idOrName string) (*models.ConcurrencyLimitV2, error) {
	var limit models.ConcurrencyLimitV2
	path := joinPath("/v2/concurrency_limits", idOrName)
	if err := s.client.get(ctx, path, &limit); err != nil {
		return nil, fmt.Errorf("failed to get v2 concurrency limit: %w", err)
	}
	return &limit, nil
}

// Update updates a v2 concurrency limit.
func (s *ConcurrencyLimitsV2Service) Update(ctx context.Context, idOrName string, req *models.ConcurrencyLimitV2Update) error {
	path := joinPath("/v2/concurrency_limits", idOrName)
	if err := s.client.patch(ctx, path, req, nil); err != nil {
		return fmt.Errorf("failed to update v2 concurrency limit: %w", err)
	}
	return nil
}

// Delete deletes a v2 concurrency limit by ID or name.
func (s *ConcurrencyLimitsV2Service) Delete(ctx context.Context, idOrName string) error {
	path := joinPath("/v2/concurrency_limits", idOrName)
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete v2 concurrency limit: %w", err)
	}
	return nil
}

// List retrieves all v2 concurrency limits.
func (s *ConcurrencyLimitsV2Service) List(ctx context.Context) ([]models.GlobalConcurrencyLimitResponse, error) {
	var limits []models.GlobalConcurrencyLimitResponse
	if err := s.client.post(ctx, "/v2/concurrency_limits/filter", nil, &limits); err != nil {
		return nil, fmt.Errorf("failed to list v2 concurrency limits: %w", err)
	}
	return limits, nil
}

// Increment increments active slots for v2 concurrency limits.
func (s *ConcurrencyLimitsV2Service) Increment(ctx context.Context, names []string, slots int, mode string) error {
	req := struct {
		Names []string `json:"names"`
		Slots int      `json:"slots"`
		Mode  string   `json:"mode,omitempty"`
	}{
		Names: names,
		Slots: slots,
		Mode:  mode,
	}

	if err := s.client.post(ctx, "/v2/concurrency_limits/increment", req, nil); err != nil {
		return fmt.Errorf("failed to increment v2 concurrency limits: %w", err)
	}
	return nil
}

// IncrementWithLease increments active slots and returns a lease.
func (s *ConcurrencyLimitsV2Service) IncrementWithLease(ctx context.Context, names []string, slots int, mode string) (*models.ConcurrencyLimitWithLeaseResponse, error) {
	req := struct {
		Names []string `json:"names"`
		Slots int      `json:"slots"`
		Mode  string   `json:"mode,omitempty"`
	}{
		Names: names,
		Slots: slots,
		Mode:  mode,
	}

	var resp models.ConcurrencyLimitWithLeaseResponse
	if err := s.client.post(ctx, "/v2/concurrency_limits/increment-with-lease", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to increment v2 concurrency limits with lease: %w", err)
	}
	return &resp, nil
}

// Decrement decrements active slots for v2 concurrency limits.
func (s *ConcurrencyLimitsV2Service) Decrement(ctx context.Context, names []string, slots int) error {
	req := struct {
		Names []string `json:"names"`
		Slots int      `json:"slots"`
	}{
		Names: names,
		Slots: slots,
	}

	if err := s.client.post(ctx, "/v2/concurrency_limits/decrement", req, nil); err != nil {
		return fmt.Errorf("failed to decrement v2 concurrency limits: %w", err)
	}
	return nil
}

// DecrementWithLease decrements active slots using a lease ID.
func (s *ConcurrencyLimitsV2Service) DecrementWithLease(ctx context.Context, names []string, slots int, leaseID uuid.UUID) error {
	req := struct {
		Names   []string  `json:"names"`
		Slots   int       `json:"slots"`
		LeaseID uuid.UUID `json:"lease_id"`
	}{
		Names:   names,
		Slots:   slots,
		LeaseID: leaseID,
	}

	if err := s.client.post(ctx, "/v2/concurrency_limits/decrement-with-lease", req, nil); err != nil {
		return fmt.Errorf("failed to decrement v2 concurrency limits with lease: %w", err)
	}
	return nil
}

// RenewLease renews a concurrency lease.
func (s *ConcurrencyLimitsV2Service) RenewLease(ctx context.Context, leaseID uuid.UUID) error {
	path := joinPath("/v2/concurrency_limits/leases", leaseID.String(), "renew")
	if err := s.client.post(ctx, path, nil, nil); err != nil {
		return fmt.Errorf("failed to renew concurrency lease: %w", err)
	}
	return nil
}
