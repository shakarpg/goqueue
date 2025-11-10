package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shakarpg/goqueue/internal/handlers"
	"github.com/shakarpg/goqueue/internal/queue"
	"go.uber.org/zap"
)

func NewRouter(queue *queue.JobQueue, logger *zap.Logger) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Handler
	jobHandler := handlers.NewJobHandler(queue, logger)

	// Rotas
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(api chi.Router) {
		api.Route("/jobs", func(jobs chi.Router) {
			jobs.Post("/", jobHandler.CreateJob)
			jobs.Get("/", jobHandler.GetAllJobs)
			jobs.Get("/{id}", jobHandler.GetJob)
		})

		api.Get("/metrics", jobHandler.GetStats)
	})

	return r
}
