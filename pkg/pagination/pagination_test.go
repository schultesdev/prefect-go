package pagination

import (
	"context"
	"errors"
	"testing"
)

func TestNewIterator(t *testing.T) {
	fetchFunc := func(ctx context.Context, offset, limit int) (*PaginatedResponse[string], error) {
		return nil, nil
	}

	tests := []struct {
		name      string
		limit     int
		wantLimit int
	}{
		{"default limit", 0, 100},
		{"custom limit", 50, 50},
		{"negative limit", -10, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iter := NewIterator(fetchFunc, tt.limit)
			if iter.limit != tt.wantLimit {
				t.Errorf("limit = %d, want %d", iter.limit, tt.wantLimit)
			}
		})
	}
}

func TestIterator_Next(t *testing.T) {
	items := []string{"item1", "item2", "item3", "item4", "item5"}

	fetchFunc := func(ctx context.Context, offset, limit int) (*PaginatedResponse[string], error) {
		if offset >= len(items) {
			return &PaginatedResponse[string]{
				Results: []string{},
				Count:   len(items),
				Limit:   limit,
				Offset:  offset,
				HasMore: false,
			}, nil
		}

		end := offset + limit
		if end > len(items) {
			end = len(items)
		}

		return &PaginatedResponse[string]{
			Results: items[offset:end],
			Count:   len(items),
			Limit:   limit,
			Offset:  offset,
			HasMore: end < len(items),
		}, nil
	}

	ctx := context.Background()
	iter := NewIterator(fetchFunc, 2)

	var results []string
	for iter.Next(ctx) {
		if item := iter.Value(); item != nil {
			results = append(results, *item)
		}
	}

	if err := iter.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(items) {
		t.Errorf("got %d items, want %d", len(results), len(items))
	}

	for i, want := range items {
		if i >= len(results) || results[i] != want {
			t.Errorf("item[%d] = %v, want %v", i, results[i], want)
		}
	}
}

func TestIterator_Error(t *testing.T) {
	expectedErr := errors.New("fetch error")

	fetchFunc := func(ctx context.Context, offset, limit int) (*PaginatedResponse[string], error) {
		return nil, expectedErr
	}

	ctx := context.Background()
	iter := NewIterator(fetchFunc, 10)

	if iter.Next(ctx) {
		t.Error("Next should return false on error")
	}

	if err := iter.Err(); err != expectedErr {
		t.Errorf("Err() = %v, want %v", err, expectedErr)
	}
}

func TestIterator_EmptyResults(t *testing.T) {
	fetchFunc := func(ctx context.Context, offset, limit int) (*PaginatedResponse[string], error) {
		return &PaginatedResponse[string]{
			Results: []string{},
			Count:   0,
			Limit:   limit,
			Offset:  offset,
			HasMore: false,
		}, nil
	}

	ctx := context.Background()
	iter := NewIterator(fetchFunc, 10)

	if iter.Next(ctx) {
		t.Error("Next should return false for empty results")
	}
}

func TestIterator_Reset(t *testing.T) {
	items := []string{"item1", "item2", "item3"}

	fetchFunc := func(ctx context.Context, offset, limit int) (*PaginatedResponse[string], error) {
		if offset >= len(items) {
			return &PaginatedResponse[string]{
				Results: []string{},
				HasMore: false,
			}, nil
		}

		return &PaginatedResponse[string]{
			Results: items[offset:],
			HasMore: false,
		}, nil
	}

	ctx := context.Background()
	iter := NewIterator(fetchFunc, 10)

	// First iteration
	count1 := 0
	for iter.Next(ctx) {
		count1++
	}

	// Reset and iterate again
	iter.Reset()
	count2 := 0
	for iter.Next(ctx) {
		count2++
	}

	if count1 != count2 {
		t.Errorf("counts don't match after reset: %d vs %d", count1, count2)
	}
	if count1 != len(items) {
		t.Errorf("got %d items, want %d", count1, len(items))
	}
}

func TestIterator_Collect(t *testing.T) {
	items := []string{"item1", "item2", "item3", "item4", "item5"}

	fetchFunc := func(ctx context.Context, offset, limit int) (*PaginatedResponse[string], error) {
		if offset >= len(items) {
			return &PaginatedResponse[string]{
				Results: []string{},
				HasMore: false,
			}, nil
		}

		end := offset + limit
		if end > len(items) {
			end = len(items)
		}

		return &PaginatedResponse[string]{
			Results: items[offset:end],
			HasMore: end < len(items),
		}, nil
	}

	ctx := context.Background()
	iter := NewIterator(fetchFunc, 2)

	results, err := iter.Collect(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(items) {
		t.Errorf("got %d items, want %d", len(results), len(items))
	}
}

func TestIterator_CollectWithLimit(t *testing.T) {
	items := []string{"item1", "item2", "item3", "item4", "item5"}

	fetchFunc := func(ctx context.Context, offset, limit int) (*PaginatedResponse[string], error) {
		if offset >= len(items) {
			return &PaginatedResponse[string]{
				Results: []string{},
				HasMore: false,
			}, nil
		}

		end := offset + limit
		if end > len(items) {
			end = len(items)
		}

		return &PaginatedResponse[string]{
			Results: items[offset:end],
			HasMore: end < len(items),
		}, nil
	}

	ctx := context.Background()
	iter := NewIterator(fetchFunc, 2)

	maxItems := 3
	results, err := iter.CollectWithLimit(ctx, maxItems)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != maxItems {
		t.Errorf("got %d items, want %d", len(results), maxItems)
	}
}

func TestPage_HasNextPage(t *testing.T) {
	tests := []struct {
		name string
		page Page[string]
		want bool
	}{
		{
			name: "has next page",
			page: Page[string]{
				Items:  []string{"a", "b"},
				Offset: 0,
				Limit:  2,
				Total:  5,
			},
			want: true,
		},
		{
			name: "no next page",
			page: Page[string]{
				Items:  []string{"d", "e"},
				Offset: 3,
				Limit:  2,
				Total:  5,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.page.HasNextPage(); got != tt.want {
				t.Errorf("HasNextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPage_HasPrevPage(t *testing.T) {
	tests := []struct {
		name string
		page Page[string]
		want bool
	}{
		{
			name: "has previous page",
			page: Page[string]{
				Offset: 2,
			},
			want: true,
		},
		{
			name: "no previous page",
			page: Page[string]{
				Offset: 0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.page.HasPrevPage(); got != tt.want {
				t.Errorf("HasPrevPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPage_NextOffset(t *testing.T) {
	page := Page[string]{
		Offset: 10,
		Limit:  5,
	}

	if got := page.NextOffset(); got != 15 {
		t.Errorf("NextOffset() = %v, want 15", got)
	}
}

func TestPage_PrevOffset(t *testing.T) {
	tests := []struct {
		name   string
		offset int
		limit  int
		want   int
	}{
		{"normal case", 10, 5, 5},
		{"at beginning", 3, 5, 0},
		{"exact boundary", 5, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := Page[string]{
				Offset: tt.offset,
				Limit:  tt.limit,
			}
			if got := page.PrevOffset(); got != tt.want {
				t.Errorf("PrevOffset() = %v, want %v", got, tt.want)
			}
		})
	}
}
