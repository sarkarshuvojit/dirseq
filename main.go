package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	configDirName = ".config/spd"
	dbFileName    = "mem.db"
	tableName     = "sequences"
)

var padding int

func main() {
	flag.IntVar(&padding, "p", 0, "Pad the output with leading zeros to reach the specified length (e.g., -p 4 => 0055)")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get user home directory", "error", err)
		os.Exit(1)
	}

	configDirPath := filepath.Join(homeDir, configDirName)

	err = os.MkdirAll(configDirPath, 0755)
	if err != nil {
		slog.Error("Failed to create config directory", "path", configDirPath, "error", err)
		os.Exit(1)
	}

	dbPath := filepath.Join(configDirPath, dbFileName)
	db, err := setupDatabase(dbPath)
	if err != nil {
		slog.Error("Failed to set up database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	currentPath, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to get current working directory", "error", err)
		os.Exit(1)
	}
	absPath, err := filepath.Abs(currentPath)
	if err != nil {
		slog.Error("Failed to get absolute path", "for", currentPath, "error", err)
		os.Exit(1)
	}

	lastSeq, err := getSequence(db, absPath)
	if err != nil {
		slog.Error("Failed to retrieve sequence", "path", absPath, "error", err)
		os.Exit(1)
	}

	newSeq := lastSeq + 1

	fmt.Printf("%0*d\n", padding, newSeq)

	err = updateSequence(db, absPath, newSeq)
	if err != nil {
		slog.Error("Failed to update sequence", "path", absPath, "new_sequence", newSeq, "error", err)
		os.Exit(1)
	}

}

// setupDatabase remains largely the same internally but returns errors
func setupDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database '%s': %w", dbPath, err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database '%s': %w", dbPath, err)
	}

	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		abs_path VARCHAR NOT NULL PRIMARY KEY,
		last_seq INT NOT NULL DEFAULT 0
	);`, tableName)

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating table '%s': %w", tableName, err)
	}

	return db, nil
}

func getSequence(db *sql.DB, path string) (int, error) {
	var lastSeq int
	querySQL := fmt.Sprintf("SELECT last_seq FROM %s WHERE abs_path = ?", tableName)

	row := db.QueryRow(querySQL, path)
	err := row.Scan(&lastSeq)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Path not found is not an application error here
		}
		return 0, fmt.Errorf("error scanning sequence row for path '%s': %w", path, err)
	}
	return lastSeq, nil
}

func updateSequence(db *sql.DB, path string, newSeq int) error {
	upsertSQL := fmt.Sprintf(`
	INSERT INTO %s (abs_path, last_seq) VALUES (?, ?)
	ON CONFLICT(abs_path) DO UPDATE SET last_seq = excluded.last_seq;
	`, tableName)

	_, err := db.Exec(upsertSQL, path, newSeq)
	if err != nil {
		return fmt.Errorf("error executing upsert for path '%s' with sequence %d: %w", path, newSeq, err)
	}
	return nil
}
