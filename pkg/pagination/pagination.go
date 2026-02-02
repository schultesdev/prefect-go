// Package pagination provides pagination support for API responses.
package pagination

import (
	"context"
	"fmt"
)

// PaginatedResponse represents a paginated response from the API.
type PaginatedResponse[T any] struct {
	// Results contains the items in this page
	Results []T

	// Count is the total number of items (if available)
	Count int

	// Limit is the maximum number of items per page
	Limit int

	// Offset is the offset of the first item in this page
	Offset int

	// HasMore indicates if there are more items to fetch
	HasMore bool
}

// FetchFunc is a function that fetches a page of results.
type FetchFunc[T any] func(ctx context.Context, offset, limit int) (*PaginatedResponse[T], error)

// Iterator provides iteration over paginated results.
type Iterator[T any] struct {
	fetchFunc FetchFunc[T]
	limit     int
	offset    int
	current   *PaginatedResponse[T]
	index     int
	err       error
	done      bool
}

// NewIterator creates a new pagination iterator.
func NewIterator[T any](fetchFunc FetchFunc[T], limit int) *Iterator[T] {
	if limit <= 0 {
		limit = 100 // Default page size
	}
	return &Iterator[T]{
		fetchFunc: fetchFunc,
		limit:     limit,
		offset:    0,
		index:     -1,
	}
}

// Next advances the iterator to the next item.
// It returns true if there is an item available, false otherwise.
func (i *Iterator[T]) Next(ctx context.Context) bool {
	if i.done || i.err != nil {
		return false
	}

	// If we don't have a current page or we've reached the end of it, fetch the next page
	if i.current == nil || i.index >= len(i.current.Results)-1 {
		// Check if we know there are no more pages
		if i.current != nil && !i.current.HasMore {
			i.done = true
			return false
		}

		// Fetch the next page
		page, err := i.fetchFunc(ctx, i.offset, i.limit)
		if err != nil {
			i.err = err
			return false
		}

		// Check if the page is empty
		if page == nil || len(page.Results) == 0 {
			i.done = true
			return false
		}

		i.current = page
		i.index = -1
		i.offset += i.limit
	}

	// Move to the next item in the current page
	i.index++
	return i.index < len(i.current.Results)
}

// Value returns the current item.
// It should only be called after Next returns true.
func (i *Iterator[T]) Value() *T {
	if i.current == nil || i.index < 0 || i.index >= len(i.current.Results) {
		return nil
	}
	return &i.current.Results[i.index]
}

// Err returns any error that occurred during iteration.
func (i *Iterator[T]) Err() error {
	return i.err
}

// Reset resets the iterator to the beginning.
func (i *Iterator[T]) Reset() {
	i.offset = 0
	i.index = -1
	i.current = nil
	i.err = nil
	i.done = false
}

// Collect fetches all items from all pages and returns them as a slice.
// Warning: This can be memory intensive for large result sets.
func (i *Iterator[T]) Collect(ctx context.Context) ([]T, error) {
	var results []T
	for i.Next(ctx) {
		if item := i.Value(); item != nil {
			results = append(results, *item)
		}
	}
	if err := i.Err(); err != nil {
		return nil, fmt.Errorf("failed to collect all items: %w", err)
	}
	return results, nil
}

// CollectWithLimit fetches items up to a maximum count.
func (i *Iterator[T]) CollectWithLimit(ctx context.Context, maxItems int) ([]T, error) {
	var results []T
	count := 0
	for i.Next(ctx) && count < maxItems {
		if item := i.Value(); item != nil {
			results = append(results, *item)
			count++
		}
	}
	if err := i.Err(); err != nil {
		return nil, fmt.Errorf("failed to collect items: %w", err)
	}
	return results, nil
}

// Page represents a single page of results with metadata.
type Page[T any] struct {
	Items  []T
	Offset int
	Limit  int
	Total  int
}

// HasNextPage returns true if there are more pages after this one.
func (p *Page[T]) HasNextPage() bool {
	return p.Offset+len(p.Items) < p.Total
}

// HasPrevPage returns true if there are pages before this one.
func (p *Page[T]) HasPrevPage() bool {
	return p.Offset > 0
}

// NextOffset returns the offset for the next page.
func (p *Page[T]) NextOffset() int {
	return p.Offset + p.Limit
}

// PrevOffset returns the offset for the previous page.
func (p *Page[T]) PrevOffset() int {
	offset := p.Offset - p.Limit
	if offset < 0 {
		return 0
	}
	return offset
}
