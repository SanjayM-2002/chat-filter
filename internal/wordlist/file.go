package wordlist

import (
	"encoding/json"
	"fmt"
	"os"
)

// File is the shape of data/words.json.
type File struct {
	Words []string `json:"words"`
}

// LoadFile reads a JSON word list from disk.
func LoadFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read word list %s: %w", path, err)
	}

	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse word list %s: %w", path, err)
	}
	return f.Words, nil
}
