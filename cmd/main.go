package main

import (
	"auto-messager/config"
	"auto-messager/internal/api"
	"auto-messager/internal/app"
	"auto-messager/internal/worker"
	"context"
	"log"
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

	ctx := context.Background()
	mainApp, err := app.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
		return
	}

	listener := worker.NewListener(mainApp)
	router := api.Router(mainApp, listener)
	mainApp.StartHTTP(router)
	listener.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mainApp.Shutdown(sctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}

	log.Println("Server exiting")
}
