package storage

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

var embeddedDB embed.FS

func resolvePath(filename string) ([]byte, error) {
	if _, err := os.Stat(filename); err == nil {
		return os.ReadFile(filename)
	}
	return embeddedDB.ReadFile(filename)
}

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
	schemaFile, err := resolvePath("db/schema.sql")
	if err != nil {
		return fmt.Errorf("failed read schema: %w", err)
	}

	if _, err := db.Exec(context.Background(), string(schemaFile)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}

func AddSeed(db *pgxpool.Pool) error {
	var count int
	err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM messages").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check messages count: %w", err)
	}

	if count > 0 {
		return nil
	}

	seedFile, err := resolvePath("db/seed.sql")
	if err != nil {
		return fmt.Errorf("failed read seed file: %w", err)
	}

	if _, err := db.Exec(context.Background(), string(seedFile)); err != nil {
		return fmt.Errorf("failed to execute seed: %w", err)
	}

	return nil
}
