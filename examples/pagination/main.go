// Pagination example demonstrating how to iterate through large result sets.
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

	// Create multiple flows for demonstration
	fmt.Println("Creating sample flows...")
	createdFlows := make([]models.Flow, 0)
	for i := 1; i <= 15; i++ {
		flow, err := c.Flows.Create(ctx, &models.FlowCreate{
			Name: fmt.Sprintf("pagination-flow-%d", i),
			Tags: []string{"pagination-example", fmt.Sprintf("batch-%d", (i-1)/5+1)},
		})
		if err != nil {
			log.Fatalf("Failed to create flow %d: %v", i, err)
		}
		createdFlows = append(createdFlows, *flow)
		fmt.Printf("✓ Created flow %d: %s\n", i, flow.Name)
	}

	// Example 1: Manual pagination with List
	fmt.Println("\n--- Example 1: Manual Pagination ---")
	pageSize := 5
	page := 0
	totalSeen := 0

	for {
		fmt.Printf("\nFetching page %d (offset: %d, limit: %d)...\n", page+1, page*pageSize, pageSize)
		result, err := c.Flows.List(ctx, &models.FlowFilter{
			Tags: []string{"pagination-example"},
		}, page*pageSize, pageSize)
		if err != nil {
			log.Fatalf("Failed to list flows: %v", err)
		}

		fmt.Printf("✓ Page %d: Found %d flows\n", page+1, len(result.Results))
		for i, flow := range result.Results {
			fmt.Printf("  %d. %s (tags: %v)\n", totalSeen+i+1, flow.Name, flow.Tags)
		}

		totalSeen += len(result.Results)

		if !result.HasMore {
			fmt.Println("\n✓ No more pages")
			break
		}

		page++
	}

	fmt.Printf("\nTotal flows seen: %d\n", totalSeen)

	// Example 2: Automatic pagination with Iterator
	fmt.Println("\n--- Example 2: Automatic Pagination with Iterator ---")
	iter := c.Flows.ListAll(ctx, &models.FlowFilter{
		Tags: []string{"pagination-example"},
	})

	count := 0
	for iter.Next(ctx) {
		flow := iter.Value()
		if flow != nil {
			count++
			fmt.Printf("%d. %s\n", count, flow.Name)
		}
	}

	if err := iter.Err(); err != nil {
		log.Fatalf("Iterator error: %v", err)
	}

	fmt.Printf("\n✓ Total flows iterated: %d\n", count)

	// Example 3: Collect all results at once
	fmt.Println("\n--- Example 3: Collect All Results ---")
	iter2 := c.Flows.ListAll(ctx, &models.FlowFilter{
		Tags: []string{"pagination-example"},
	})

	allFlows, err := iter2.Collect(ctx)
	if err != nil {
		log.Fatalf("Failed to collect flows: %v", err)
	}

	fmt.Printf("✓ Collected %d flows at once\n", len(allFlows))
	for i, flow := range allFlows {
		fmt.Printf("  %d. %s\n", i+1, flow.Name)
	}

	// Example 4: Collect with limit
	fmt.Println("\n--- Example 4: Collect With Limit ---")
	iter3 := c.Flows.ListAll(ctx, &models.FlowFilter{
		Tags: []string{"pagination-example"},
	})

	maxItems := 7
	limitedFlows, err := iter3.CollectWithLimit(ctx, maxItems)
	if err != nil {
		log.Fatalf("Failed to collect flows: %v", err)
	}

	fmt.Printf("✓ Collected first %d flows (limit: %d)\n", len(limitedFlows), maxItems)
	for i, flow := range limitedFlows {
		fmt.Printf("  %d. %s\n", i+1, flow.Name)
	}

	// Example 5: Filter by specific tag batch
	fmt.Println("\n--- Example 5: Filter by Specific Tag ---")
	batch2Filter := &models.FlowFilter{
		Tags: []string{"batch-2"},
	}

	batch2Page, err := c.Flows.List(ctx, batch2Filter, 0, 10)
	if err != nil {
		log.Fatalf("Failed to list flows: %v", err)
	}

	fmt.Printf("✓ Found %d flows with tag 'batch-2'\n", len(batch2Page.Results))
	for i, flow := range batch2Page.Results {
		fmt.Printf("  %d. %s (tags: %v)\n", i+1, flow.Name, flow.Tags)
	}

	// Clean up all created flows
	fmt.Println("\n--- Cleaning Up ---")
	for i, flow := range createdFlows {
		if err := c.Flows.Delete(ctx, flow.ID); err != nil {
			log.Printf("Warning: Failed to delete flow %d: %v", i+1, err)
		} else {
			fmt.Printf("✓ Deleted flow %d: %s\n", i+1, flow.Name)
		}
	}

	fmt.Println("\n✓ Pagination example completed successfully!")
}
