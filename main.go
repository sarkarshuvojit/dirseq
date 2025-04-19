package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	dirseq "github.com/sarkarshuvojit/dirseq/pkg"
)

const (
	configDirName = ".config/direseq"
)

var padding int


func setupFlags() {
	flag.IntVar(&padding, "p", 0, "Pad the output with leading zeros to reach the specified length (e.g., -p 4 => 0055)")
	flag.Parse()
}

func main() {
	setupFlags()
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

	db, err := dirseq.SetupDatabase(configDirPath)
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

	lastSeq, err := dirseq.GetSequence(db, absPath)
	if err != nil {
		slog.Error("Failed to retrieve sequence", "path", absPath, "error", err)
		os.Exit(1)
	}

	newSeq := lastSeq + 1

	fmt.Printf("%0*d\n", padding, newSeq)

	err = dirseq.UpdateSequence(db, absPath, newSeq)
	if err != nil {
		slog.Error("Failed to update sequence", "path", absPath, "new_sequence", newSeq, "error", err)
		os.Exit(1)
	}

}
