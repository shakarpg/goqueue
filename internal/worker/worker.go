package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/shakarpg/goqueue/internal/models"
	"github.com/shakarpg/goqueue/internal/queue"
	"go.uber.org/zap"
)

type WorkerPool struct {
	numWorkers int
	queue      *queue.JobQueue
	logger     *zap.Logger
}

func NewWorkerPool(numWorkers int, queue *queue.JobQueue, logger *zap.Logger) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		queue:      queue,
		logger:     logger,
	}
}

// Start inicia os workers
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 1; i <= wp.numWorkers; i++ {
		go wp.worker(ctx, i)
	}
	wp.logger.Info("✅ Workers iniciados", zap.Int("count", wp.numWorkers))
}

// worker processa jobs da fila
func (wp *WorkerPool) worker(ctx context.Context, id int) {
	wp.logger.Info("🔧 Worker iniciado", zap.Int("worker_id", id))

	for {
		select {
		case <-ctx.Done():
			wp.logger.Info("🛑 Worker finalizado", zap.Int("worker_id", id))
			return
		case job := <-wp.queue.Dequeue():
			wp.processJob(job, id)
		}
	}
}

// processJob processa um job específico
func (wp *WorkerPool) processJob(job *models.Job, workerID int) {
	now := time.Now()
	job.Status = models.StatusRunning
	job.StartedAt = &now
	wp.queue.UpdateJob(job)

	wp.logger.Info("⚙️  Processando job",
		zap.String("job_id", job.ID),
		zap.String("type", string(job.Type)),
		zap.Int("worker_id", workerID),
	)

	// Simular processamento baseado no tipo
	var err error
	var result string

	switch job.Type {
	case models.JobTypeEmail:
		result, err = wp.processEmail(job)
	case models.JobTypePDF:
		result, err = wp.processPDF(job)
	case models.JobTypeImage:
		result, err = wp.processImage(job)
	default:
		err = fmt.Errorf("tipo de job desconhecido: %s", job.Type)
	}

	endTime := time.Now()
	job.EndedAt = &endTime

	if err != nil {
		job.Status = models.StatusFailed
		job.Error = err.Error()
		wp.logger.Error("❌ Job falhou",
			zap.String("job_id", job.ID),
			zap.Error(err),
		)
	} else {
		job.Status = models.StatusCompleted
		job.Result = result
		wp.logger.Info("✅ Job completado",
			zap.String("job_id", job.ID),
			zap.Duration("duration", endTime.Sub(*job.StartedAt)),
		)
	}

	wp.queue.UpdateJob(job)
}

// processEmail simula envio de email
func (wp *WorkerPool) processEmail(job *models.Job) (string, error) {
	time.Sleep(2 * time.Second) // Simular processamento
	to := job.Payload["to"]
	return fmt.Sprintf("Email enviado para %v", to), nil
}

// processPDF simula geração de PDF
func (wp *WorkerPool) processPDF(job *models.Job) (string, error) {
	time.Sleep(3 * time.Second) // Simular processamento
	filename := job.Payload["filename"]
	return fmt.Sprintf("PDF gerado: %v", filename), nil
}

// processImage simula processamento de imagem
func (wp *WorkerPool) processImage(job *models.Job) (string, error) {
	time.Sleep(4 * time.Second) // Simular processamento
	url := job.Payload["url"]
	return fmt.Sprintf("Imagem processada: %v", url), nil
}
