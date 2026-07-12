package db

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate() {
	// Create the tracking table first — always safe to run
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name       TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatalf("could not create schema_migrations table: %v", err)
	}

	// fs.ReadDir returns entries in lexical order — 001, 002, 003 etc.
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		log.Fatalf("could not read embedded migrations: %v", err)
	}

	for _, entry := range entries {
		name := entry.Name()

		// Check whether this migration has already been applied
		var count int
		DB.QueryRow(
			"SELECT COUNT(*) FROM schema_migrations WHERE name = ?", name,
		).Scan(&count)

		if count > 0 {
			continue // already applied — skip silently
		}

		sql, err := migrationFiles.ReadFile("migrations/" + name)
		if err != nil {
			log.Fatalf("could not read migration %s: %v", name, err)
		}

		// Wrap each migration in a transaction so it is all-or-nothing
		tx, err := DB.Begin()
		if err != nil {
			log.Fatalf("migration %s: could not begin transaction: %v", name, err)
		}

		if _, err := tx.Exec(string(sql)); err != nil {
			tx.Rollback()
			log.Fatalf("migration %s failed, rolled back: %v", name, err)
		}

		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (name) VALUES (?)", name,
		); err != nil {
			tx.Rollback()
			log.Fatalf("migration %s: could not record in schema_migrations: %v", name, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("migration %s: could not commit: %v", name, err)
		}

		log.Printf("migration applied: %s", name)
	}
}
