package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"assignment_2_AP/internal/api"
	"assignment_2_AP/internal/queue"
	"assignment_2_AP/internal/store"
	"assignment_2_AP/internal/worker"
)

func main() {
	taskStore := store.NewRepository[string, *store.Task]()
	taskQueue := queue.NewQueue[*store.Task](100)

	stats := store.NewStats()
	// 3 work with pool
	workerPool := worker.NewPool(3, taskQueue, taskStore, stats)

	// Start worker pool
	workerPool.Start()

	// background monitoring
	monitor := worker.NewMonitor(taskStore, 5*time.Second)
	monitor.Start()

	// HTTP handlers
	handler := api.NewHandler(taskStore, taskQueue, stats)

	http.HandleFunc("/tasks", handler.TasksHandler)
	http.HandleFunc("/tasks/", handler.TaskByIDHandler)
	http.HandleFunc("/stats", handler.StatsHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	// start  goroutine for server
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	monitor.Stop()
	workerPool.Stop()

	log.Println("server stopped")
}
