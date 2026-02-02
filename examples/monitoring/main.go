// Monitoring example demonstrating how to monitor flow runs until completion.
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

	// Create a flow
	fmt.Println("Creating flow...")
	flow, err := c.Flows.Create(ctx, &models.FlowCreate{
		Name: "monitored-flow",
		Tags: []string{"monitoring-example"},
	})
	if err != nil {
		log.Fatalf("Failed to create flow: %v", err)
	}
	fmt.Printf("✓ Flow created: %s (ID: %s)\n", flow.Name, flow.ID)

	// Create a flow run
	fmt.Println("\nCreating flow run...")
	flowRun, err := c.FlowRuns.Create(ctx, &models.FlowRunCreate{
		FlowID: flow.ID,
		Name:   "monitored-run",
	})
	if err != nil {
		log.Fatalf("Failed to create flow run: %v", err)
	}
	fmt.Printf("✓ Flow run created: %s (ID: %s)\n", flowRun.Name, flowRun.ID)

	// Start monitoring in a goroutine
	go func() {
		// Simulate state transitions
		time.Sleep(2 * time.Second)

		// Set to RUNNING
		fmt.Println("\n[Simulator] Setting state to RUNNING...")
		runningState := models.StateTypeRunning
		_, err := c.FlowRuns.SetState(ctx, flowRun.ID, &models.StateCreate{
			Type:    runningState,
			Message: strPtr("Flow run started"),
		})
		if err != nil {
			log.Printf("Error setting state: %v", err)
		}

		// Simulate work
		time.Sleep(5 * time.Second)

		// Set to COMPLETED
		fmt.Println("[Simulator] Setting state to COMPLETED...")
		completedState := models.StateTypeCompleted
		_, err = c.FlowRuns.SetState(ctx, flowRun.ID, &models.StateCreate{
			Type:    completedState,
			Message: strPtr("Flow run completed successfully"),
		})
		if err != nil {
			log.Printf("Error setting state: %v", err)
		}
	}()

	// Monitor the flow run until completion
	fmt.Println("\nMonitoring flow run (polling every 2 seconds)...")
	fmt.Println("Waiting for completion...")

	// Create a context with timeout
	monitorCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Wait for the flow run to complete
	finalRun, err := c.FlowRuns.Wait(monitorCtx, flowRun.ID, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to wait for flow run: %v", err)
	}

	fmt.Printf("\n✓ Flow run completed!\n")
	fmt.Printf("  Final state: %v\n", finalRun.StateType)
	if finalRun.State != nil && finalRun.State.Message != nil {
		fmt.Printf("  Message: %s\n", *finalRun.State.Message)
	}
	if finalRun.StartTime != nil {
		fmt.Printf("  Start time: %v\n", finalRun.StartTime)
	}
	if finalRun.EndTime != nil {
		fmt.Printf("  End time: %v\n", finalRun.EndTime)
	}
	fmt.Printf("  Total run time: %.2f seconds\n", finalRun.TotalRunTime)

	// Query flow runs with a filter
	fmt.Println("\nQuerying completed flow runs...")
	completedState := models.StateTypeCompleted
	runsPage, err := c.FlowRuns.List(ctx, &models.FlowRunFilter{
		FlowID:    &flow.ID,
		StateType: &completedState,
	}, 0, 10)
	if err != nil {
		log.Fatalf("Failed to list flow runs: %v", err)
	}
	fmt.Printf("✓ Found %d completed flow runs\n", len(runsPage.Results))

	// Clean up
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

	fmt.Println("\n✓ Monitoring example completed successfully!")
}

func strPtr(s string) *string {
	return &s
}
