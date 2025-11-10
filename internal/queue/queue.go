package queue

import (
	"sync"

	"github.com/shakarpg/goqueue/internal/models"
)

type JobQueue struct {
	jobs    chan *models.Job
	storage map[string]*models.Job
	mu      sync.RWMutex
}

func NewJobQueue(size int) *JobQueue {
	return &JobQueue{
		jobs:    make(chan *models.Job, size),
		storage: make(map[string]*models.Job),
	}
}

// Enqueue adiciona um job na fila
func (q *JobQueue) Enqueue(job *models.Job) {
	q.mu.Lock()
	q.storage[job.ID] = job
	q.mu.Unlock()

	q.jobs <- job
}

// Dequeue retorna o próximo job da fila
func (q *JobQueue) Dequeue() <-chan *models.Job {
	return q.jobs
}

// GetJob retorna um job pelo ID
func (q *JobQueue) GetJob(id string) (*models.Job, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	job, exists := q.storage[id]
	return job, exists
}

// GetAllJobs retorna todos os jobs
func (q *JobQueue) GetAllJobs() []*models.Job {
	q.mu.RLock()
	defer q.mu.RUnlock()

	jobs := make([]*models.Job, 0, len(q.storage))
	for _, job := range q.storage {
		jobs = append(jobs, job)
	}
	return jobs
}

// UpdateJob atualiza um job no storage
func (q *JobQueue) UpdateJob(job *models.Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.storage[job.ID] = job
}

// GetStats retorna estatísticas da fila
func (q *JobQueue) GetStats() map[string]int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	stats := map[string]int{
		"total":     len(q.storage),
		"pending":   0,
		"running":   0,
		"completed": 0,
		"failed":    0,
	}

	for _, job := range q.storage {
		switch job.Status {
		case models.StatusPending:
			stats["pending"]++
		case models.StatusRunning:
			stats["running"]++
		case models.StatusCompleted:
			stats["completed"]++
		case models.StatusFailed:
			stats["failed"]++
		}
	}

	return stats
}
