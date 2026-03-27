package client

import (
	"context"
	"fmt"
	"net/http"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// AutomationsService handles operations related to automations.
type AutomationsService struct {
	client *Client
}

// Create creates a new automation.
func (s *AutomationsService) Create(ctx context.Context, req *models.AutomationCreate) (*models.Automation, error) {
	var automation models.Automation
	if err := s.client.do(ctx, http.MethodPost, "/automations/", req, &automation); err != nil {
		return nil, fmt.Errorf("failed to create automation: %w", err)
	}
	return &automation, nil
}

// Get retrieves an automation by ID.
func (s *AutomationsService) Get(ctx context.Context, id uuid.UUID) (*models.Automation, error) {
	var automation models.Automation
	path := joinPath("/automations", id.String())
	if err := s.client.get(ctx, path, &automation); err != nil {
		return nil, fmt.Errorf("failed to get automation: %w", err)
	}
	return &automation, nil
}

// Update fully replaces an automation.
func (s *AutomationsService) Update(ctx context.Context, id uuid.UUID, req *models.AutomationUpdate) (*models.Automation, error) {
	var automation models.Automation
	path := joinPath("/automations", id.String())
	if err := s.client.do(ctx, http.MethodPut, path, req, &automation); err != nil {
		return nil, fmt.Errorf("failed to update automation: %w", err)
	}
	return &automation, nil
}

// Patch partially updates an automation.
func (s *AutomationsService) Patch(ctx context.Context, id uuid.UUID, req *models.AutomationPartialUpdate) error {
	path := joinPath("/automations", id.String())
	if err := s.client.patch(ctx, path, req, nil); err != nil {
		return fmt.Errorf("failed to patch automation: %w", err)
	}
	return nil
}

// Delete deletes an automation by ID.
func (s *AutomationsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/automations", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete automation: %w", err)
	}
	return nil
}

// List retrieves automations with optional filtering.
func (s *AutomationsService) List(ctx context.Context, filter *models.AutomationFilter) ([]models.Automation, error) {
	if filter == nil {
		filter = &models.AutomationFilter{}
	}

	var automations []models.Automation
	if err := s.client.post(ctx, "/automations/filter", filter, &automations); err != nil {
		return nil, fmt.Errorf("failed to list automations: %w", err)
	}
	return automations, nil
}

// Count returns the number of automations matching the filter.
func (s *AutomationsService) Count(ctx context.Context, filter *models.AutomationFilter) (int, error) {
	if filter == nil {
		filter = &models.AutomationFilter{}
	}

	var count int
	if err := s.client.post(ctx, "/automations/count", filter, &count); err != nil {
		return 0, fmt.Errorf("failed to count automations: %w", err)
	}
	return count, nil
}

// GetRelatedTo retrieves automations related to a specific resource.
func (s *AutomationsService) GetRelatedTo(ctx context.Context, resourceID string) ([]models.Automation, error) {
	var automations []models.Automation
	path := joinPath("/automations/related-to", resourceID)
	if err := s.client.get(ctx, path, &automations); err != nil {
		return nil, fmt.Errorf("failed to get related automations: %w", err)
	}
	return automations, nil
}

// DeleteOwnedBy deletes automations owned by a specific resource.
func (s *AutomationsService) DeleteOwnedBy(ctx context.Context, resourceID string) error {
	path := joinPath("/automations/owned-by", resourceID)
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete owned automations: %w", err)
	}
	return nil
}

// ValidateTemplate validates an automation template.
func (s *AutomationsService) ValidateTemplate(ctx context.Context, template string) error {
	req := struct {
		Template string `json:"template"`
	}{
		Template: template,
	}
	if err := s.client.post(ctx, "/automations/templates/validate", req, nil); err != nil {
		return fmt.Errorf("failed to validate template: %w", err)
	}
	return nil
}
