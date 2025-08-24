package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

func InitPostgre(dbUri string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dbUri)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return db, nil
}

// TODO: dynamic schema path @kucukahmet
func ExecSchema(db *pgxpool.Pool) error {
	schemaFile, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return fmt.Errorf("failed read schema: %w", err)
	}

	if _, err := db.Exec(context.Background(), string(schemaFile)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}
