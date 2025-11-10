package tests

import (
	"testing"
	"time"

	"github.com/shakarpg/goqueue/internal/models"
	"github.com/shakarpg/goqueue/internal/queue"
	"github.com/stretchr/testify/assert"
)

func TestJobQueue(t *testing.T) {
	q := queue.NewJobQueue(10)

	job := &models.Job{
		ID:        "test-1",
		Type:      models.JobTypeEmail,
		Status:    models.StatusPending,
		Payload:   map[string]interface{}{"to": "test@example.com"},
		CreatedAt: time.Now(),
	}

	// Enqueue
	q.Enqueue(job)

	// GetJob
	retrieved, exists := q.GetJob("test-1")
	assert.True(t, exists)
	assert.Equal(t, job.ID, retrieved.ID)

	// Stats
	stats := q.GetStats()
	assert.Equal(t, 1, stats["total"])
	assert.Equal(t, 1, stats["pending"])
}

func TestJobQueueStats(t *testing.T) {
	q := queue.NewJobQueue(10)

	jobs := []*models.Job{
		{ID: "1", Status: models.StatusPending, CreatedAt: time.Now()},
		{ID: "2", Status: models.StatusRunning, CreatedAt: time.Now()},
		{ID: "3", Status: models.StatusCompleted, CreatedAt: time.Now()},
		{ID: "4", Status: models.StatusFailed, CreatedAt: time.Now()},
	}

	for _, job := range jobs {
		q.Enqueue(job)
	}

	stats := q.GetStats()
	assert.Equal(t, 4, stats["total"])
	assert.Equal(t, 1, stats["pending"])
	assert.Equal(t, 1, stats["running"])
	assert.Equal(t, 1, stats["completed"])
	assert.Equal(t, 1, stats["failed"])
}
