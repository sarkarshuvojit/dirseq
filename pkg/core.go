package dirseq

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func Show(
	configDirName string,
	padding int,
) {
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

	db, err := SetupDatabase(configDirPath)
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

	lastSeq, err := GetSequence(db, absPath)
	if err != nil {
		slog.Error("Failed to retrieve sequence", "path", absPath, "error", err)
		os.Exit(1)
	}

	newSeq := lastSeq + 1

	fmt.Printf("%0*d\n", padding, newSeq)

	err = UpdateSequence(db, absPath, newSeq)
	if err != nil {
		slog.Error("Failed to update sequence", "path", absPath, "new_sequence", newSeq, "error", err)
		os.Exit(1)
	}
}
