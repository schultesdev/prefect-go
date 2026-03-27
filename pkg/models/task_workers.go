package models

import "time"

// TaskWorkerFilter represents filter criteria for querying task workers.
type TaskWorkerFilter struct {
	TaskKeys []string `json:"task_keys"`
}

// TaskWorkerResponse represents a task worker response.
type TaskWorkerResponse struct {
	Identifier string    `json:"identifier"`
	TaskKeys   []string  `json:"task_keys"`
	Timestamp  time.Time `json:"timestamp"`
}
