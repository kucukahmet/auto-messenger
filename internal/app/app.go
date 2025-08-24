package app

import (
	"auto-messager/config"
	"auto-messager/internal/storage"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config  *config.Config
	DB      *pgxpool.Pool
	Queries *storage.Queries
	HTTP    *http.Server
}

func NewApp(ctx context.Context, config *config.Config) (*App, error) {
	db, err := storage.InitPostgre(config.DB_URI)
	if err != nil {
		return nil, err
	}
	queries := storage.New(db)

	return &App{
		Config:  config,
		DB:      db,
		Queries: queries,
	}, nil
}

func (app *App) StartHTTP(handler http.Handler) {
	app.HTTP = &http.Server{
		Addr:         ":" + app.Config.HTTP_PORT,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Printf("Server starting on port %s", app.Config.HTTP_PORT)
		if err := app.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()
}

func (app *App) Shutdown(ctx context.Context) error {
	if app.HTTP != nil {
		_ = app.HTTP.Shutdown(ctx)
	}

	if app.DB != nil {
		app.DB.Close()
	}
	return nil
}
