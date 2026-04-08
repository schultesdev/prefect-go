package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"git.schultes.dev/schultesdev/prefect-go/pkg/pagination"
	"github.com/google/uuid"
)

// DeploymentsService handles operations related to deployments.
type DeploymentsService struct {
	client *Client
}

// Create creates a new deployment.
func (s *DeploymentsService) Create(ctx context.Context, req *models.DeploymentCreate) (*models.Deployment, error) {
	var deployment models.Deployment
	if err := s.client.post(ctx, "/deployments/", req, &deployment); err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}
	return &deployment, nil
}

// Get retrieves a deployment by ID.
func (s *DeploymentsService) Get(ctx context.Context, id uuid.UUID) (*models.Deployment, error) {
	var deployment models.Deployment
	path := joinPath("/deployments", id.String())
	if err := s.client.get(ctx, path, &deployment); err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return &deployment, nil
}

// GetByName retrieves a deployment by flow name and deployment name.
func (s *DeploymentsService) GetByName(ctx context.Context, flowName, deploymentName string) (*models.Deployment, error) {
	var deployment models.Deployment
	path := fmt.Sprintf("/deployments/name/%s/%s", flowName, deploymentName)
	if err := s.client.get(ctx, path, &deployment); err != nil {
		return nil, fmt.Errorf("failed to get deployment by name: %w", err)
	}
	return &deployment, nil
}

// Update updates a deployment.
func (s *DeploymentsService) Update(ctx context.Context, id uuid.UUID, req *models.DeploymentUpdate) (*models.Deployment, error) {
	var deployment models.Deployment
	path := joinPath("/deployments", id.String())
	if err := s.client.patch(ctx, path, req, &deployment); err != nil {
		return nil, fmt.Errorf("failed to update deployment: %w", err)
	}
	return &deployment, nil
}

// Delete deletes a deployment by ID.
func (s *DeploymentsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/deployments", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}
	return nil
}

// Pause pauses a deployment's schedule.
func (s *DeploymentsService) Pause(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/deployments", id.String(), "pause")
	if err := s.client.post(ctx, path, nil, nil); err != nil {
		return fmt.Errorf("failed to pause deployment: %w", err)
	}
	return nil
}

// Resume resumes a deployment's schedule.
func (s *DeploymentsService) Resume(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/deployments", id.String(), "resume")
	if err := s.client.post(ctx, path, nil, nil); err != nil {
		return fmt.Errorf("failed to resume deployment: %w", err)
	}
	return nil
}

// CreateFlowRun creates a flow run from a deployment.
func (s *DeploymentsService) CreateFlowRun(ctx context.Context, id uuid.UUID, params map[string]interface{}) (*models.FlowRun, error) {
	var flowRun models.FlowRun
	path := joinPath("/deployments", id.String(), "create_flow_run")

	req := struct {
		Parameters map[string]interface{} `json:"parameters,omitempty"`
	}{
		Parameters: params,
	}

	if err := s.client.post(ctx, path, req, &flowRun); err != nil {
		return nil, fmt.Errorf("failed to create flow run from deployment: %w", err)
	}
	return &flowRun, nil
}

// List retrieves a list of deployments with optional filtering.
func (s *DeploymentsService) List(ctx context.Context, filter *models.DeploymentFilter, offset, limit int) (*pagination.PaginatedResponse[models.Deployment], error) {
	if filter == nil {
		filter = &models.DeploymentFilter{}
	}
	filter.Offset = offset
	filter.Limit = limit

	type response struct {
		Results []models.Deployment `json:"results"`
		Count   int                 `json:"count"`
	}

	var resp response
	if err := s.client.post(ctx, "/deployments/filter", filter, &resp); err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	return &pagination.PaginatedResponse[models.Deployment]{
		Results: resp.Results,
		Count:   resp.Count,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(resp.Results) < resp.Count,
	}, nil
}

// ListAll returns an iterator for all deployments matching the filter.
func (s *DeploymentsService) ListAll(ctx context.Context, filter *models.DeploymentFilter) *pagination.Iterator[models.Deployment] {
	fetchFunc := func(ctx context.Context, offset, limit int) (*pagination.PaginatedResponse[models.Deployment], error) {
		return s.List(ctx, filter, offset, limit)
	}
	return pagination.NewIterator(fetchFunc, 100)
}

// Count returns the number of deployments matching the filter.
func (s *DeploymentsService) Count(ctx context.Context, filter *models.DeploymentFilter) (int, error) {
	if filter == nil {
		filter = &models.DeploymentFilter{}
	}

	var count int
	if err := s.client.post(ctx, "/deployments/count", filter, &count); err != nil {
		return 0, fmt.Errorf("failed to count deployments: %w", err)
	}
	return count, nil
}

// GetSchedules retrieves all schedules for a deployment.
func (s *DeploymentsService) GetSchedules(ctx context.Context, id uuid.UUID) ([]models.DeploymentSchedule, error) {
	var schedules []models.DeploymentSchedule
	path := joinPath("/deployments", id.String(), "schedules")
	if err := s.client.get(ctx, path, &schedules); err != nil {
		return nil, fmt.Errorf("failed to get deployment schedules: %w", err)
	}
	return schedules, nil
}

// CreateSchedules creates schedules for a deployment.
func (s *DeploymentsService) CreateSchedules(ctx context.Context, id uuid.UUID, schedules []models.DeploymentScheduleCreate) ([]models.DeploymentSchedule, error) {
	var result []models.DeploymentSchedule
	path := joinPath("/deployments", id.String(), "schedules")
	if err := s.client.post(ctx, path, schedules, &result); err != nil {
		return nil, fmt.Errorf("failed to create deployment schedules: %w", err)
	}
	return result, nil
}

// UpdateSchedule updates a specific schedule for a deployment.
func (s *DeploymentsService) UpdateSchedule(ctx context.Context, deploymentID, scheduleID uuid.UUID, req *models.DeploymentScheduleUpdate) error {
	path := joinPath("/deployments", deploymentID.String(), "schedules", scheduleID.String())
	if err := s.client.patch(ctx, path, req, nil); err != nil {
		return fmt.Errorf("failed to update deployment schedule: %w", err)
	}
	return nil
}

// DeleteSchedule deletes a specific schedule from a deployment.
func (s *DeploymentsService) DeleteSchedule(ctx context.Context, deploymentID, scheduleID uuid.UUID) error {
	path := joinPath("/deployments", deploymentID.String(), "schedules", scheduleID.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete deployment schedule: %w", err)
	}
	return nil
}
