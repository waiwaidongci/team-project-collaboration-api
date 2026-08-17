package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"team-project-task-api/internal/config"
)

func main() {
	databaseURL := flag.String("database-url", "", "PostgreSQL connection URL")
	migrationsPath := flag.String("path", "migrations", "directory containing .sql migrations")
	flag.Parse()

	cfg := config.Load()
	if *databaseURL == "" {
		*databaseURL = cfg.DatabaseURL
	}

	db, err := sql.Open("postgres", *databaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`); err != nil {
		log.Fatalf("create schema_migrations: %v", err)
	}

	files, err := migrationFiles(*migrationsPath)
	if err != nil {
		log.Fatalf("list migrations: %v", err)
	}

	for _, file := range files {
		version := filepath.Base(file)
		applied, err := migrationApplied(ctx, db, version)
		if err != nil {
			log.Fatalf("check migration %s: %v", version, err)
		}
		if applied {
			log.Printf("skip %s", version)
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("read migration %s: %v", file, err)
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Fatalf("begin migration %s: %v", version, err)
		}
		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			log.Fatalf("apply migration %s: %v", version, err)
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			_ = tx.Rollback()
			log.Fatalf("record migration %s: %v", version, err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("commit migration %s: %v", version, err)
		}
		log.Printf("applied %s", version)
	}
	fmt.Println("migrations complete")
}

func migrationFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(files)
	return files, nil
}

func migrationApplied(ctx context.Context, db *sql.DB, version string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
	return exists, err
}
