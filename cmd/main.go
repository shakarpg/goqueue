package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shakarpg/goqueue/internal/queue"
	"github.com/shakarpg/goqueue/internal/router"
	"github.com/shakarpg/goqueue/internal/worker"
	"go.uber.org/zap"
)

func main() {
	// Configurar logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Criar fila de jobs
	jobQueue := queue.NewJobQueue(100)

	// Iniciar workers (5 workers concorrentes)
	numWorkers := 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := worker.NewWorkerPool(numWorkers, jobQueue, logger)
	pool.Start(ctx)

	// Configurar rotas
	r := router.NewRouter(jobQueue, logger)

	// Configurar servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Canal para capturar sinais de shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		logger.Info("🚀 Servidor rodando", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("❌ Erro ao iniciar servidor", zap.Error(err))
		}
	}()

	// Aguardar sinal de shutdown
	<-quit
	logger.Info("🛑 Desligando servidor...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("❌ Erro no shutdown", zap.Error(err))
	}

	// Cancelar workers
	cancel()
	time.Sleep(2 * time.Second) // Aguardar workers finalizarem

	logger.Info("✅ Servidor desligado com sucesso")
}
