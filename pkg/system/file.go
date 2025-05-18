package system

import (
	"fmt"
	"os"
	"path/filepath"
)

// StoreFile stores a file on the file system at the given path.
// If the path does not exist, it will be created.
func StoreFile(path string, content []byte) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write the file
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
