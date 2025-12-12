package mv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qawatake/hiden/internal/mkdir"
)

// Run moves a file to the date-based hiden directory in the current git repository.
// It creates the directory if it doesn't exist.
func Run(dirname string, filePath string) (string, error) {
	// Ensure the target directory exists
	targetDir, relDir, err := mkdir.EnsureDir(dirname)
	if err != nil {
		return "", err
	}

	// Get the base name of the file
	baseName := filepath.Base(filePath)

	// Construct the target file path
	targetPath := filepath.Join(targetDir, baseName)

	// Move the file
	if err := os.Rename(filePath, targetPath); err != nil {
		return "", fmt.Errorf("failed to move file: %w", err)
	}

	// Return path relative to repository root
	relPath := filepath.Join(relDir, baseName)
	return relPath, nil
}
