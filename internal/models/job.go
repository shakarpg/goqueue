package models

import "time"

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type JobType string

const (
	JobTypeEmail JobType = "email"
	JobTypePDF   JobType = "pdf"
	JobTypeImage JobType = "image"
)

type Job struct {
	ID        string                 `json:"id"`
	Type      JobType                `json:"type"`
	Status    JobStatus              `json:"status"`
	Payload   map[string]interface{} `json:"payload"`
	Result    string                 `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	StartedAt *time.Time             `json:"started_at,omitempty"`
	EndedAt   *time.Time             `json:"ended_at,omitempty"`
}
