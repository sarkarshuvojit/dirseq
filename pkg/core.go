package dirseq

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func Show(
	padding int,
) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get user home directory", "error", err)
		os.Exit(1)
	}

	configDirPath := filepath.Join(homeDir, ConfigDirPath)

	err = os.MkdirAll(configDirPath, 0700)
	if err != nil {
		slog.Error("Failed to create config directory", "path", configDirPath, "error", err)
		os.Exit(1)
	}
	store := &JsonStore{}
	_, err = store.SetupDatabase()
	if err != nil {
		slog.Error("Failed to set up database", "error", err)
		os.Exit(1)
	}

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

	lastSeq, configPadding, err := store.GetSequenceAndPadding(absPath)
	if err != nil {
		slog.Error("Failed to retrieve sequence", "path", absPath, "error", err)
		os.Exit(1)
	}

	newSeq := lastSeq + 1
	if padding > 0 {
		configPadding = padding
	}
	fmt.Printf("%0*d\n", configPadding, newSeq)

	err = store.UpdateSequence(absPath, newSeq)
	if err != nil {
		slog.Error("Failed to update sequence", "path", absPath, "new_sequence", newSeq, "error", err)
		os.Exit(1)
	}
}
