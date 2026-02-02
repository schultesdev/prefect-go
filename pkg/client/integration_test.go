//go:build integration
// +build integration

package client_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gregor/prefect-go/pkg/client"
	"github.com/gregor/prefect-go/pkg/errors"
	"github.com/gregor/prefect-go/pkg/models"
)

var testClient *client.Client

func TestMain(m *testing.M) {
	// Get base URL from environment or use default
	baseURL := os.Getenv("PREFECT_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:4200/api"
	}

	// Create test client
	var err error
	testClient, err = client.NewClient(
		client.WithBaseURL(baseURL),
	)
	if err != nil {
		panic(err)
	}

	// Run tests
	os.Exit(m.Run())
}

func TestIntegration_FlowLifecycle(t *testing.T) {
	ctx := context.Background()

	// Create a flow
	flow, err := testClient.Flows.Create(ctx, &models.FlowCreate{
		Name: "test-flow-" + uuid.New().String(),
		Tags: []string{"integration-test"},
	})
	if err != nil {
		t.Fatalf("Failed to create flow: %v", err)
	}
	defer testClient.Flows.Delete(ctx, flow.ID)

	// Verify flow was created
	if flow.ID == uuid.Nil {
		t.Error("Flow ID is nil")
	}
	if flow.Name == "" {
		t.Error("Flow name is empty")
	}

	// Get the flow
	retrievedFlow, err := testClient.Flows.Get(ctx, flow.ID)
	if err != nil {
		t.Fatalf("Failed to get flow: %v", err)
	}
	if retrievedFlow.ID != flow.ID {
		t.Errorf("Flow ID mismatch: got %v, want %v", retrievedFlow.ID, flow.ID)
	}

	// Update the flow
	newTags := []string{"integration-test", "updated"}
	updatedFlow, err := testClient.Flows.Update(ctx, flow.ID, &models.FlowUpdate{
		Tags: &newTags,
	})
	if err != nil {
		t.Fatalf("Failed to update flow: %v", err)
	}
	if len(updatedFlow.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(updatedFlow.Tags))
	}

	// List flows
	flowsPage, err := testClient.Flows.List(ctx, &models.FlowFilter{
		Tags: []string{"integration-test"},
	}, 0, 10)
	if err != nil {
		t.Fatalf("Failed to list flows: %v", err)
	}
	if len(flowsPage.Results) == 0 {
		t.Error("Expected at least one flow in results")
	}

	// Delete the flow
	if err := testClient.Flows.Delete(ctx, flow.ID); err != nil {
		t.Fatalf("Failed to delete flow: %v", err)
	}

	// Verify deletion
	_, err = testClient.Flows.Get(ctx, flow.ID)
	if !errors.IsNotFound(err) {
		t.Errorf("Expected NotFound error, got: %v", err)
	}
}

func TestIntegration_FlowRunLifecycle(t *testing.T) {
	ctx := context.Background()

	// Create a flow first
	flow, err := testClient.Flows.Create(ctx, &models.FlowCreate{
		Name: "test-flow-run-" + uuid.New().String(),
		Tags: []string{"integration-test"},
	})
	if err != nil {
		t.Fatalf("Failed to create flow: %v", err)
	}
	defer testClient.Flows.Delete(ctx, flow.ID)

	// Create a flow run
	flowRun, err := testClient.FlowRuns.Create(ctx, &models.FlowRunCreate{
		FlowID: flow.ID,
		Name:   "test-run",
		Parameters: map[string]interface{}{
			"test_param": "value",
		},
	})
	if err != nil {
		t.Fatalf("Failed to create flow run: %v", err)
	}
	defer testClient.FlowRuns.Delete(ctx, flowRun.ID)

	// Verify flow run was created
	if flowRun.ID == uuid.Nil {
		t.Error("Flow run ID is nil")
	}
	if flowRun.FlowID != flow.ID {
		t.Errorf("Flow ID mismatch: got %v, want %v", flowRun.FlowID, flow.ID)
	}

	// Get the flow run
	retrievedRun, err := testClient.FlowRuns.Get(ctx, flowRun.ID)
	if err != nil {
		t.Fatalf("Failed to get flow run: %v", err)
	}
	if retrievedRun.ID != flowRun.ID {
		t.Errorf("Flow run ID mismatch: got %v, want %v", retrievedRun.ID, flowRun.ID)
	}

	// Set state to RUNNING
	runningState := models.StateTypeRunning
	updatedRun, err := testClient.FlowRuns.SetState(ctx, flowRun.ID, &models.StateCreate{
		Type:    runningState,
		Message: strPtr("Test running"),
	})
	if err != nil {
		t.Fatalf("Failed to set state: %v", err)
	}
	if updatedRun.StateType == nil || *updatedRun.StateType != runningState {
		t.Errorf("State type mismatch: got %v, want %v", updatedRun.StateType, runningState)
	}

	// Set state to COMPLETED
	completedState := models.StateTypeCompleted
	completedRun, err := testClient.FlowRuns.SetState(ctx, flowRun.ID, &models.StateCreate{
		Type:    completedState,
		Message: strPtr("Test completed"),
	})
	if err != nil {
		t.Fatalf("Failed to set state: %v", err)
	}
	if completedRun.StateType == nil || *completedRun.StateType != completedState {
		t.Errorf("State type mismatch: got %v, want %v", completedRun.StateType, completedState)
	}

	// List flow runs
	runsPage, err := testClient.FlowRuns.List(ctx, &models.FlowRunFilter{
		FlowID: &flow.ID,
	}, 0, 10)
	if err != nil {
		t.Fatalf("Failed to list flow runs: %v", err)
	}
	if len(runsPage.Results) == 0 {
		t.Error("Expected at least one flow run in results")
	}
}

func TestIntegration_DeploymentLifecycle(t *testing.T) {
	ctx := context.Background()

	// Create a flow first
	flow, err := testClient.Flows.Create(ctx, &models.FlowCreate{
		Name: "test-deployment-flow-" + uuid.New().String(),
		Tags: []string{"integration-test"},
	})
	if err != nil {
		t.Fatalf("Failed to create flow: %v", err)
	}
	defer testClient.Flows.Delete(ctx, flow.ID)

	// Create a deployment
	workPoolName := "default-pool"
	deployment, err := testClient.Deployments.Create(ctx, &models.DeploymentCreate{
		Name:         "test-deployment",
		FlowID:       flow.ID,
		WorkPoolName: &workPoolName,
	})
	if err != nil {
		t.Fatalf("Failed to create deployment: %v", err)
	}
	defer testClient.Deployments.Delete(ctx, deployment.ID)

	// Verify deployment was created
	if deployment.ID == uuid.Nil {
		t.Error("Deployment ID is nil")
	}
	if deployment.FlowID != flow.ID {
		t.Errorf("Flow ID mismatch: got %v, want %v", deployment.FlowID, flow.ID)
	}

	// Get the deployment
	retrievedDeployment, err := testClient.Deployments.Get(ctx, deployment.ID)
	if err != nil {
		t.Fatalf("Failed to get deployment: %v", err)
	}
	if retrievedDeployment.ID != deployment.ID {
		t.Errorf("Deployment ID mismatch: got %v, want %v", retrievedDeployment.ID, deployment.ID)
	}

	// Pause the deployment
	if err := testClient.Deployments.Pause(ctx, deployment.ID); err != nil {
		t.Fatalf("Failed to pause deployment: %v", err)
	}

	// Verify paused
	pausedDeployment, err := testClient.Deployments.Get(ctx, deployment.ID)
	if err != nil {
		t.Fatalf("Failed to get deployment: %v", err)
	}
	if !pausedDeployment.Paused {
		t.Error("Expected deployment to be paused")
	}

	// Resume the deployment
	if err := testClient.Deployments.Resume(ctx, deployment.ID); err != nil {
		t.Fatalf("Failed to resume deployment: %v", err)
	}

	// Verify resumed
	resumedDeployment, err := testClient.Deployments.Get(ctx, deployment.ID)
	if err != nil {
		t.Fatalf("Failed to get deployment: %v", err)
	}
	if resumedDeployment.Paused {
		t.Error("Expected deployment to be resumed")
	}
}

func TestIntegration_Pagination(t *testing.T) {
	ctx := context.Background()

	// Create multiple flows
	createdFlows := make([]uuid.UUID, 0)
	for i := 0; i < 15; i++ {
		flow, err := testClient.Flows.Create(ctx, &models.FlowCreate{
			Name: "pagination-test-" + uuid.New().String(),
			Tags: []string{"pagination-integration-test"},
		})
		if err != nil {
			t.Fatalf("Failed to create flow %d: %v", i, err)
		}
		createdFlows = append(createdFlows, flow.ID)
	}

	// Clean up
	defer func() {
		for _, id := range createdFlows {
			testClient.Flows.Delete(ctx, id)
		}
	}()

	// Test manual pagination
	page1, err := testClient.Flows.List(ctx, &models.FlowFilter{
		Tags: []string{"pagination-integration-test"},
	}, 0, 5)
	if err != nil {
		t.Fatalf("Failed to list flows: %v", err)
	}
	if len(page1.Results) != 5 {
		t.Errorf("Expected 5 flows in page 1, got %d", len(page1.Results))
	}
	if !page1.HasMore {
		t.Error("Expected more pages")
	}

	// Test iterator
	iter := testClient.Flows.ListAll(ctx, &models.FlowFilter{
		Tags: []string{"pagination-integration-test"},
	})

	count := 0
	for iter.Next(ctx) {
		count++
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("Iterator error: %v", err)
	}
	if count != 15 {
		t.Errorf("Expected to iterate over 15 flows, got %d", count)
	}
}

func TestIntegration_Variables(t *testing.T) {
	ctx := context.Background()

	// Create a variable
	variable, err := testClient.Variables.Create(ctx, &models.VariableCreate{
		Name:  "test-var-" + uuid.New().String(),
		Value: "test-value",
		Tags:  []string{"integration-test"},
	})
	if err != nil {
		t.Fatalf("Failed to create variable: %v", err)
	}
	defer testClient.Variables.Delete(ctx, variable.ID)

	// Get the variable
	retrievedVar, err := testClient.Variables.Get(ctx, variable.ID)
	if err != nil {
		t.Fatalf("Failed to get variable: %v", err)
	}
	if retrievedVar.Value != "test-value" {
		t.Errorf("Value mismatch: got %v, want test-value", retrievedVar.Value)
	}

	// Update the variable
	newValue := "updated-value"
	updatedVar, err := testClient.Variables.Update(ctx, variable.ID, &models.VariableUpdate{
		Value: &newValue,
	})
	if err != nil {
		t.Fatalf("Failed to update variable: %v", err)
	}
	if updatedVar.Value != newValue {
		t.Errorf("Value mismatch: got %v, want %v", updatedVar.Value, newValue)
	}
}

func TestIntegration_AdminEndpoints(t *testing.T) {
	ctx := context.Background()

	// Test health check
	if err := testClient.Admin.Health(ctx); err != nil {
		t.Fatalf("Health check failed: %v", err)
	}

	// Test version
	version, err := testClient.Admin.Version(ctx)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}
	if version == "" {
		t.Error("Version is empty")
	}
	t.Logf("Server version: %s", version)
}

func TestIntegration_FlowRunWait(t *testing.T) {
	ctx := context.Background()

	// Create a flow
	flow, err := testClient.Flows.Create(ctx, &models.FlowCreate{
		Name: "wait-test-" + uuid.New().String(),
		Tags: []string{"integration-test"},
	})
	if err != nil {
		t.Fatalf("Failed to create flow: %v", err)
	}
	defer testClient.Flows.Delete(ctx, flow.ID)

	// Create a flow run
	flowRun, err := testClient.FlowRuns.Create(ctx, &models.FlowRunCreate{
		FlowID: flow.ID,
		Name:   "wait-test-run",
	})
	if err != nil {
		t.Fatalf("Failed to create flow run: %v", err)
	}
	defer testClient.FlowRuns.Delete(ctx, flowRun.ID)

	// Simulate completion in background
	go func() {
		time.Sleep(2 * time.Second)
		completedState := models.StateTypeCompleted
		testClient.FlowRuns.SetState(ctx, flowRun.ID, &models.StateCreate{
			Type:    completedState,
			Message: strPtr("Test completed"),
		})
	}()

	// Wait for completion with timeout
	waitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	finalRun, err := testClient.FlowRuns.Wait(waitCtx, flowRun.ID, time.Second)
	if err != nil {
		t.Fatalf("Failed to wait for flow run: %v", err)
	}

	if finalRun.StateType == nil || *finalRun.StateType != models.StateTypeCompleted {
		t.Errorf("Expected COMPLETED state, got %v", finalRun.StateType)
	}
}

func strPtr(s string) *string {
	return &s
}
