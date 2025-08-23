package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitPostgre(dbUri string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		_ = db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func ExecSchema(db *sql.DB) error {
	schemaFile, err := os.ReadFile("db/schema.sql")
	if err != nil {
		return fmt.Errorf("failed read schema: %w", err)
	}
	if _, err := db.Exec(string(schemaFile)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}
