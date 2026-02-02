// Basic example demonstrating how to create and run flows with the Prefect Go client.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gregor/prefect-go/pkg/client"
	"github.com/gregor/prefect-go/pkg/models"
)

func main() {
	// Create a new Prefect client
	c, err := client.NewClient(
		client.WithBaseURL("http://localhost:4200/api"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Create a new flow
	fmt.Println("Creating flow...")
	flow, err := c.Flows.Create(ctx, &models.FlowCreate{
		Name: "example-flow",
		Tags: []string{"example", "go-client", "basic"},
		Labels: map[string]interface{}{
			"environment": "development",
			"version":     "1.0",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create flow: %v", err)
	}
	fmt.Printf("✓ Flow created: %s (ID: %s)\n", flow.Name, flow.ID)

	// Get the flow by ID
	fmt.Println("\nRetrieving flow...")
	retrievedFlow, err := c.Flows.Get(ctx, flow.ID)
	if err != nil {
		log.Fatalf("Failed to get flow: %v", err)
	}
	fmt.Printf("✓ Flow retrieved: %s\n", retrievedFlow.Name)

	// Create a flow run
	fmt.Println("\nCreating flow run...")
	flowRun, err := c.FlowRuns.Create(ctx, &models.FlowRunCreate{
		FlowID: flow.ID,
		Name:   "example-run-1",
		Parameters: map[string]interface{}{
			"param1": "value1",
			"param2": 42,
		},
		Tags: []string{"manual-run"},
	})
	if err != nil {
		log.Fatalf("Failed to create flow run: %v", err)
	}
	fmt.Printf("✓ Flow run created: %s (ID: %s)\n", flowRun.Name, flowRun.ID)
	fmt.Printf("  State: %v\n", flowRun.StateType)

	// Set the flow run state to RUNNING
	fmt.Println("\nUpdating flow run state to RUNNING...")
	runningState := models.StateTypeRunning
	updatedRun, err := c.FlowRuns.SetState(ctx, flowRun.ID, &models.StateCreate{
		Type:    runningState,
		Message: strPtr("Flow run started"),
	})
	if err != nil {
		log.Fatalf("Failed to set state: %v", err)
	}
	fmt.Printf("✓ State updated to: %v\n", updatedRun.StateType)

	// Simulate some work
	fmt.Println("\nSimulating work...")
	time.Sleep(2 * time.Second)

	// Set the flow run state to COMPLETED
	fmt.Println("Updating flow run state to COMPLETED...")
	completedState := models.StateTypeCompleted
	completedRun, err := c.FlowRuns.SetState(ctx, flowRun.ID, &models.StateCreate{
		Type:    completedState,
		Message: strPtr("Flow run completed successfully"),
	})
	if err != nil {
		log.Fatalf("Failed to set state: %v", err)
	}
	fmt.Printf("✓ State updated to: %v\n", completedRun.StateType)

	// List all flows
	fmt.Println("\nListing all flows...")
	flowsPage, err := c.Flows.List(ctx, nil, 0, 10)
	if err != nil {
		log.Fatalf("Failed to list flows: %v", err)
	}
	fmt.Printf("✓ Found %d flows (total: %d)\n", len(flowsPage.Results), flowsPage.Count)
	for i, f := range flowsPage.Results {
		fmt.Printf("  %d. %s (ID: %s)\n", i+1, f.Name, f.ID)
	}

	// List flow runs for this flow
	fmt.Println("\nListing flow runs...")
	runsPage, err := c.FlowRuns.List(ctx, &models.FlowRunFilter{
		FlowID: &flow.ID,
	}, 0, 10)
	if err != nil {
		log.Fatalf("Failed to list flow runs: %v", err)
	}
	fmt.Printf("✓ Found %d flow runs\n", len(runsPage.Results))
	for i, run := range runsPage.Results {
		fmt.Printf("  %d. %s - State: %v\n", i+1, run.Name, run.StateType)
	}

	// Clean up: delete the flow run and flow
	fmt.Println("\nCleaning up...")
	if err := c.FlowRuns.Delete(ctx, flowRun.ID); err != nil {
		log.Printf("Warning: Failed to delete flow run: %v", err)
	} else {
		fmt.Println("✓ Flow run deleted")
	}

	if err := c.Flows.Delete(ctx, flow.ID); err != nil {
		log.Printf("Warning: Failed to delete flow: %v", err)
	} else {
		fmt.Println("✓ Flow deleted")
	}

	fmt.Println("\n✓ Example completed successfully!")
}

func strPtr(s string) *string {
	return &s
}
