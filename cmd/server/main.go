package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"taskmanager/internal/handler"
	"taskmanager/internal/repository"
	"taskmanager/internal/service"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	repo := repository.NewInMemoryTaskRepository(logger)
	svc := service.NewTaskService(repo, logger)
	taskHandler := handler.NewTaskHandler(svc, logger)
	serviceHandler := handler.NewServiceHandler()
	healthHandler := handler.NewHealthHandler()

	mux := http.NewServeMux()
	serviceHandler.RegisterRoutes(mux)
	healthHandler.RegisterRoutes(mux)
	taskHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// HTTP/2 is enabled by default for TLS servers in Go's stdlib.
		// For plaintext, Go 1.6+ supports h2c via third-party, but for now we use HTTP/1.1 for local dev.
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		logger.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown", zap.Error(err))
		}
	}()

	logger.Info("Starting server", zap.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("ListenAndServe failed", zap.Error(err))
	}
	logger.Info("Server exited cleanly")
}
