package dirseq

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type sequenceData struct {
	LastSeq int `json:"last_seq"`
	Padding int `json:"padding"`
}

type JsonStore struct {
	data     map[string]sequenceData
	filePath string
	mu       sync.Mutex
}

func (j *JsonStore) SetupDatabase() (*struct{}, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	storePath := filepath.Join(homeDir, ConfigDirPath)
	if err := os.MkdirAll(storePath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	j.filePath = filepath.Join(storePath, "sequence_store.json")
	j.data = make(map[string]sequenceData)

	if _, err := os.Stat(j.filePath); os.IsNotExist(err) {
		// File doesn't exist -> create empty file
		if err := j.persist(); err != nil {
			return nil, fmt.Errorf("failed to create empty JSON store: %w", err)
		}
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat JSON file: %w", err)
	}

	// Load existing data
	file, err := os.Open(j.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&j.data); err != nil {
		return nil, fmt.Errorf("failed to decode existing JSON: %w", err)
	}

	return nil, nil
}

func (j *JsonStore) persist() error {
	file, err := os.Create(j.filePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(j.data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

func (j *JsonStore) GetSequence(path string) (int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	data, ok := j.data[path]
	if !ok {
		return 0, nil
	}
	return data.LastSeq, nil
}

func (j *JsonStore) GetSequenceAndPadding(path string) (int, int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	data, ok := j.data[path]
	if !ok {
		return 0, 0, nil
	}
	return data.LastSeq, data.Padding, nil
}

func (j *JsonStore) UpdateSequence(path string, newSeq int) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	data := j.data[path] // if not present, this is zero-value struct
	data.LastSeq = newSeq
	j.data[path] = data
	return j.persist()
}

func (j *JsonStore) UpdatePadding(path string, padding int) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	data := j.data[path] // if not present, this is zero-value struct
	data.Padding = padding
	j.data[path] = data
	return j.persist()
}

var _ Store[struct{}] = (*JsonStore)(nil)
