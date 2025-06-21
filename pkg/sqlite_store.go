package dirseq

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type SqliteStore struct {
	db *sql.DB
}

// GetSequence implements Store.
func (s *SqliteStore) GetSequence(path string) (int, error) {
	var lastSeq int
	querySQL := fmt.Sprintf("SELECT last_seq FROM %s WHERE abs_path = ?", TableName)

	row := s.db.QueryRow(querySQL, path)
	err := row.Scan(&lastSeq)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Path not found is not an application error here
		}
		return 0, fmt.Errorf("error scanning sequence row for path '%s': %w", path, err)
	}
	return lastSeq, nil
}

// GetSequenceAndPadding implements Store.
func (s *SqliteStore) GetSequenceAndPadding(path string) (int, int, error) {
	var lastSeq int
	var padding int

	querySQL := fmt.Sprintf("SELECT last_seq, padding FROM %s WHERE abs_path = ?", TableName)

	row := s.db.QueryRow(querySQL, path)
	err := row.Scan(&lastSeq, &padding)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil // No entry for this path; return defaults
		}
		return 0, 0, fmt.Errorf("error scanning sequence row for path '%s': %w", path, err)
	}

	return lastSeq, padding, nil
}

// SetupDatabase implements Store.
func (s *SqliteStore) SetupDatabase() (*sql.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get user home directory", "error", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(homeDir, ConfigDirPath)
	dbPath = filepath.Join(dbPath, DbFileName)

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
	CREATE TABLE IF NOT EXISTS %[1]s (
		abs_path VARCHAR NOT NULL PRIMARY KEY,
		last_seq INT NOT NULL DEFAULT 0,
		padding INT NOT NULL DEFAULT 0
	);
	`, TableName)

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating table '%s': %w", TableName, err)
	}

	return db, nil
}

// UpdatePadding implements Store.
func (s *SqliteStore) UpdatePadding(path string, padding int) error {
	updatePaddingSql := fmt.Sprintf(`
		INSERT INTO %s (abs_path, padding)
		VALUES (?, ?)
		ON CONFLICT(abs_path)
		DO UPDATE SET padding = excluded.padding
	`, TableName)

	_, err := s.db.Exec(updatePaddingSql, path, padding)
	if err != nil {
		return fmt.Errorf("error executing upsert for path '%s' with sequence %d: %w", path, padding, err)
	}
	return nil
}

// UpdateSequence implements Store.
func (s *SqliteStore) UpdateSequence(path string, newSeq int) error {
	upsertSQL := fmt.Sprintf(`
	INSERT INTO %s (abs_path, last_seq) VALUES (?, ?)
	ON CONFLICT(abs_path) DO UPDATE SET last_seq = excluded.last_seq;
	`, TableName)

	_, err := s.db.Exec(upsertSQL, path, newSeq)
	if err != nil {
		return fmt.Errorf("error executing upsert for path '%s' with sequence %d: %w", path, newSeq, err)
	}
	return nil
}

var _ Store[sql.DB] = (*SqliteStore)(nil)
