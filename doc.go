/*
Package prefect provides a Go client for interacting with the Prefect Server API.

The client supports all major Prefect operations including flows, flow runs,
deployments, tasks, work pools, and more. It includes built-in retry logic,
pagination support, and comprehensive error handling.

# Installation

	go get git.schultes.dev/schultesdev/prefect-go

# Quick Start

Create a client and interact with the Prefect API:

	package main

	import (
		"context"
		"log"

		"git.schultes.dev/schultesdev/prefect-go/pkg/client"
		"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	)

	func main() {
		// Create a new client
		c, err := client.NewClient(
			client.WithBaseURL("http://localhost:4200/api"),
		)
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()

		// Create a flow
		flow, err := c.Flows.Create(ctx, &models.FlowCreate{
			Name: "my-flow",
			Tags: []string{"example"},
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Created flow: %s (ID: %s)", flow.Name, flow.ID)
	}

# Client Configuration

The client can be configured with various options:

	client, err := client.NewClient(
		client.WithBaseURL("http://localhost:4200/api"),
		client.WithAPIKey("your-api-key"),
		client.WithTimeout(30 * time.Second),
		client.WithRetry(retry.Config{
			MaxAttempts: 5,
			InitialDelay: time.Second,
		}),
	)

# Services

The client provides several services for different resource types:

  - Flows: Flow management operations
  - FlowRuns: Flow run operations and monitoring
  - Deployments: Deployment creation and management
  - TaskRuns: Task run operations
  - WorkPools: Work pool management
  - WorkQueues: Work queue operations
  - Variables: Variable storage and retrieval
  - Logs: Log creation and querying
  - Admin: Administrative operations (health, version)

# Error Handling

The package provides structured error types for API errors:

	flow, err := client.Flows.Get(ctx, flowID)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Println("Flow not found")
		} else if errors.IsUnauthorized(err) {
			log.Println("Authentication required")
		} else {
			log.Printf("API error: %v", err)
		}
	}

# Pagination

For endpoints that return lists, use pagination:

	// Manual pagination
	page, err := client.Flows.List(ctx, filter, 0, 100)

	// Automatic pagination with iterator
	iter := client.Flows.ListAll(ctx, filter)
	for iter.Next(ctx) {
		flow := iter.Value()
		// Process flow
	}
	if err := iter.Err(); err != nil {
		log.Fatal(err)
	}

# Retry Logic

The client automatically retries failed requests with exponential backoff:

  - Retries on transient errors (5xx, 429)
  - Respects Retry-After headers
  - Context-aware cancellation
  - Configurable retry behavior

# Monitoring Flow Runs

Wait for a flow run to complete:

	finalRun, err := client.FlowRuns.Wait(ctx, runID, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Flow run completed with state: %v", finalRun.StateType)

For more examples, see the examples/ directory in the repository.
*/
package prefect
