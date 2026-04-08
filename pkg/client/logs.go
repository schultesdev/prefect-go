package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"git.schultes.dev/schultesdev/prefect-go/pkg/pagination"
)

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
