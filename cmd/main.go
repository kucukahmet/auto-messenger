package main

import (
	"auto-messager/config"
	"auto-messager/internal/api"
	"auto-messager/internal/storage"
	"auto-messager/internal/worker"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	listener := worker.NewListener(int32(config.PERIOD))
	db, err := storage.InitPostgre(config.DB_URI)

	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
		return
	}

	err = storage.ExecSchema(db)
	if err != nil {
		log.Fatalf("Failed to execute schema: %v", err)
		return
	}

	router := api.Router(listener, db)

	server := &http.Server{
		Addr:         ":" + config.HTTP_PORT,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", config.HTTP_PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")
}
