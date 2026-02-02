// Deployment example demonstrating how to create and manage deployments.
package main

import (
	"context"
	"fmt"
	"log"

	"git.schultes.dev/schultesdev/prefect-go/pkg/client"
	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
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

	// Create a flow first
	fmt.Println("Creating flow...")
	flow, err := c.Flows.Create(ctx, &models.FlowCreate{
		Name: "scheduled-flow",
		Tags: []string{"scheduled", "production"},
	})
	if err != nil {
		log.Fatalf("Failed to create flow: %v", err)
	}
	fmt.Printf("✓ Flow created: %s (ID: %s)\n", flow.Name, flow.ID)

	// Create a deployment
	fmt.Println("\nCreating deployment...")
	workPoolName := "default-pool"
	deployment, err := c.Deployments.Create(ctx, &models.DeploymentCreate{
		Name:         "my-deployment",
		FlowID:       flow.ID,
		WorkPoolName: &workPoolName,
		Parameters: map[string]interface{}{
			"default_param": "value",
		},
		Tags: []string{"automated"},
	})
	if err != nil {
		log.Fatalf("Failed to create deployment: %v", err)
	}
	fmt.Printf("✓ Deployment created: %s (ID: %s)\n", deployment.Name, deployment.ID)
	fmt.Printf("  Status: %v\n", deployment.Status)
	fmt.Printf("  Paused: %v\n", deployment.Paused)

	// Get the deployment
	fmt.Println("\nRetrieving deployment...")
	retrievedDeployment, err := c.Deployments.Get(ctx, deployment.ID)
	if err != nil {
		log.Fatalf("Failed to get deployment: %v", err)
	}
	fmt.Printf("✓ Deployment retrieved: %s\n", retrievedDeployment.Name)

	// Get deployment by name
	fmt.Println("\nRetrieving deployment by name...")
	deploymentByName, err := c.Deployments.GetByName(ctx, flow.Name, deployment.Name)
	if err != nil {
		log.Fatalf("Failed to get deployment by name: %v", err)
	}
	fmt.Printf("✓ Deployment retrieved by name: %s\n", deploymentByName.Name)

	// Pause the deployment
	fmt.Println("\nPausing deployment...")
	if err := c.Deployments.Pause(ctx, deployment.ID); err != nil {
		log.Fatalf("Failed to pause deployment: %v", err)
	}
	fmt.Println("✓ Deployment paused")

	// Check paused status
	pausedDeployment, err := c.Deployments.Get(ctx, deployment.ID)
	if err != nil {
		log.Fatalf("Failed to get deployment: %v", err)
	}
	fmt.Printf("  Paused: %v\n", pausedDeployment.Paused)

	// Resume the deployment
	fmt.Println("\nResuming deployment...")
	if err := c.Deployments.Resume(ctx, deployment.ID); err != nil {
		log.Fatalf("Failed to resume deployment: %v", err)
	}
	fmt.Println("✓ Deployment resumed")

	// Create a flow run from the deployment
	fmt.Println("\nCreating flow run from deployment...")
	flowRun, err := c.Deployments.CreateFlowRun(ctx, deployment.ID, map[string]interface{}{
		"custom_param": "custom_value",
	})
	if err != nil {
		log.Fatalf("Failed to create flow run: %v", err)
	}
	fmt.Printf("✓ Flow run created: %s (ID: %s)\n", flowRun.Name, flowRun.ID)
	fmt.Printf("  Deployment ID: %v\n", flowRun.DeploymentID)

	// Update deployment description
	fmt.Println("\nUpdating deployment...")
	description := "Updated deployment description"
	updatedDeployment, err := c.Deployments.Update(ctx, deployment.ID, &models.DeploymentUpdate{
		Description: &description,
	})
	if err != nil {
		log.Fatalf("Failed to update deployment: %v", err)
	}
	fmt.Printf("✓ Deployment updated\n")
	fmt.Printf("  Description: %v\n", updatedDeployment.Description)

	// Clean up
	fmt.Println("\nCleaning up...")
	if err := c.FlowRuns.Delete(ctx, flowRun.ID); err != nil {
		log.Printf("Warning: Failed to delete flow run: %v", err)
	} else {
		fmt.Println("✓ Flow run deleted")
	}

	if err := c.Deployments.Delete(ctx, deployment.ID); err != nil {
		log.Printf("Warning: Failed to delete deployment: %v", err)
	} else {
		fmt.Println("✓ Deployment deleted")
	}

	if err := c.Flows.Delete(ctx, flow.ID); err != nil {
		log.Printf("Warning: Failed to delete flow: %v", err)
	} else {
		fmt.Println("✓ Flow deleted")
	}

	fmt.Println("\n✓ Deployment example completed successfully!")
}
