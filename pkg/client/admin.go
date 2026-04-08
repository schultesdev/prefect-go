package client

import (
	"context"
	"fmt"
)

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
