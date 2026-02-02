package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gregor/prefect-go/pkg/models"
	"github.com/gregor/prefect-go/pkg/pagination"
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

// TaskRunsService handles operations related to task runs.
type TaskRunsService struct {
	client *Client
}

// Create creates a new task run.
func (t *TaskRunsService) Create(ctx context.Context, req *models.TaskRunCreate) (*models.TaskRun, error) {
	var taskRun models.TaskRun
	if err := t.client.post(ctx, "/task_runs/", req, &taskRun); err != nil {
		return nil, fmt.Errorf("failed to create task run: %w", err)
	}
	return &taskRun, nil
}

// Get retrieves a task run by ID.
func (t *TaskRunsService) Get(ctx context.Context, id uuid.UUID) (*models.TaskRun, error) {
	var taskRun models.TaskRun
	path := joinPath("/task_runs", id.String())
	if err := t.client.get(ctx, path, &taskRun); err != nil {
		return nil, fmt.Errorf("failed to get task run: %w", err)
	}
	return &taskRun, nil
}

// Delete deletes a task run by ID.
func (t *TaskRunsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/task_runs", id.String())
	if err := t.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete task run: %w", err)
	}
	return nil
}

// SetState sets the state of a task run.
func (t *TaskRunsService) SetState(ctx context.Context, id uuid.UUID, state *models.StateCreate) (*models.TaskRun, error) {
	var taskRun models.TaskRun
	path := joinPath("/task_runs", id.String(), "set_state")

	req := struct {
		State *models.StateCreate `json:"state"`
	}{
		State: state,
	}

	if err := t.client.post(ctx, path, req, &taskRun); err != nil {
		return nil, fmt.Errorf("failed to set task run state: %w", err)
	}
	return &taskRun, nil
}

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

// LogsService handles operations related to logs.
type LogsService struct {
	client *Client
}

// Create creates new log entries.
func (l *LogsService) Create(ctx context.Context, logs []*models.LogCreate) error {
	if err := l.client.post(ctx, "/logs/", logs, nil); err != nil {
		return fmt.Errorf("failed to create logs: %w", err)
	}
	return nil
}

// List retrieves logs with filtering.
func (l *LogsService) List(ctx context.Context, filter interface{}, offset, limit int) (*pagination.PaginatedResponse[models.Log], error) {
	type request struct {
		Filter interface{} `json:"filter,omitempty"`
		Offset int         `json:"offset"`
		Limit  int         `json:"limit"`
	}

	req := request{
		Filter: filter,
		Offset: offset,
		Limit:  limit,
	}

	type response struct {
		Results []models.Log `json:"results"`
		Count   int          `json:"count"`
	}

	var resp response
	if err := l.client.post(ctx, "/logs/filter", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}

	return &pagination.PaginatedResponse[models.Log]{
		Results: resp.Results,
		Count:   resp.Count,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(resp.Results) < resp.Count,
	}, nil
}

// AdminService handles administrative operations.
type AdminService struct {
	client *Client
}

// Health checks the health of the Prefect server.
func (a *AdminService) Health(ctx context.Context) error {
	if err := a.client.get(ctx, "/health", nil); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	return nil
}

// Version retrieves the server version.
func (a *AdminService) Version(ctx context.Context) (string, error) {
	var result struct {
		Version string `json:"version"`
	}
	if err := a.client.get(ctx, "/version", &result); err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}
	return result.Version, nil
}
