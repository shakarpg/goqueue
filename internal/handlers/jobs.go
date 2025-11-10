package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shakarpg/goqueue/internal/models"
	"github.com/shakarpg/goqueue/internal/queue"
	"go.uber.org/zap"
)

type JobHandler struct {
	queue  *queue.JobQueue
	logger *zap.Logger
}

func NewJobHandler(queue *queue.JobQueue, logger *zap.Logger) *JobHandler {
	return &JobHandler{
		queue:  queue,
		logger: logger,
	}
}

type CreateJobRequest struct {
	Type    models.JobType         `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// CreateJob cria um novo job
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validar tipo de job
	if req.Type != models.JobTypeEmail && req.Type != models.JobTypePDF && req.Type != models.JobTypeImage {
		http.Error(w, "Invalid job type", http.StatusBadRequest)
		return
	}

	// Criar job
	job := &models.Job{
		ID:        uuid.New().String(),
		Type:      req.Type,
		Status:    models.StatusPending,
		Payload:   req.Payload,
		CreatedAt: time.Now(),
	}

	// Adicionar na fila
	h.queue.Enqueue(job)

	h.logger.Info("📥 Job criado", zap.String("job_id", job.ID), zap.String("type", string(job.Type)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

// GetJob retorna um job específico
func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	job, exists := h.queue.GetJob(id)
	if !exists {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// GetAllJobs retorna todos os jobs
func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs := h.queue.GetAllJobs()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// GetStats retorna estatísticas da fila
func (h *JobHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := h.queue.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
