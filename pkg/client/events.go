package client

import (
	"context"
	"fmt"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
)

// EventsService handles operations related to events.
type EventsService struct {
	client *Client
}

// Create creates new events.
func (s *EventsService) Create(ctx context.Context, events []models.Event) error {
	if err := s.client.post(ctx, "/events", events, nil); err != nil {
		return fmt.Errorf("failed to create events: %w", err)
	}
	return nil
}

// List retrieves events with filtering.
func (s *EventsService) List(ctx context.Context, filter *models.EventFilter) (*models.EventPage, error) {
	if filter == nil {
		filter = &models.EventFilter{}
	}

	var page models.EventPage
	if err := s.client.post(ctx, "/events/filter", filter, &page); err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	return &page, nil
}

// NextPage retrieves the next page of events using a page token.
func (s *EventsService) NextPage(ctx context.Context) (*models.EventPage, error) {
	var page models.EventPage
	if err := s.client.get(ctx, "/events/filter/next", &page); err != nil {
		return nil, fmt.Errorf("failed to get next events page: %w", err)
	}
	return &page, nil
}

// CountBy counts events grouped by a countable field.
func (s *EventsService) CountBy(ctx context.Context, countable string, filter *models.EventFilter) ([]models.EventCount, error) {
	if filter == nil {
		filter = &models.EventFilter{}
	}

	path := joinPath("/events/count-by", countable)
	var counts []models.EventCount
	if err := s.client.post(ctx, path, filter, &counts); err != nil {
		return nil, fmt.Errorf("failed to count events: %w", err)
	}
	return counts, nil
}
