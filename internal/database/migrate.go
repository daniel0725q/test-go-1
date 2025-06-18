package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// RunMigrations runs all pending database migrations
func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files by version
	var migrations []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrations = append(migrations, file.Name())
		}
	}
	sort.Strings(migrations)

	// Get applied migrations
	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		version := getVersionFromFilename(migration)
		if _, exists := applied[version]; exists {
			continue
		}

		// Read and execute migration
		content, err := ioutil.ReadFile(filepath.Join("migrations", migration))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migration, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration, err)
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO schema_migrations (version, dirty) VALUES ($1, false)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration, err)
		}

		log.Printf("Applied migration: %s", migration)
	}

	log.Println("Migrations completed successfully")
	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint NOT NULL,
			dirty boolean NOT NULL,
			PRIMARY KEY (version)
		)
	`)
	return err
}

func getAppliedMigrations(db *sql.DB) (map[int64]struct{}, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int64]struct{})
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = struct{}{}
	}
	return applied, rows.Err()
}

func getVersionFromFilename(filename string) int64 {
	parts := strings.Split(filename, "_")
	if len(parts) > 0 {
		version, _ := strconv.ParseInt(parts[0], 10, 64)
		return version
	}
	return 0
}
