package mv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qawatake/hiden/internal/mkdir"
)

// Run moves files to the date-based hiden directory in the current git repository.
// It creates the directory if it doesn't exist.
func Run(dirname string, filePaths []string) ([]string, error) {
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("no files specified")
	}

	// Ensure the target directory exists
	targetDir, relDir, err := mkdir.EnsureDir(dirname)
	if err != nil {
		return nil, err
	}

	var relPaths []string
	for _, filePath := range filePaths {
		// Get the base name of the file
		baseName := filepath.Base(filePath)

		// Construct the target file path
		targetPath := filepath.Join(targetDir, baseName)

		// Move the file
		if err := os.Rename(filePath, targetPath); err != nil {
			return relPaths, fmt.Errorf("failed to move file %s: %w", filePath, err)
		}

		// Collect path relative to repository root
		relPath := filepath.Join(relDir, baseName)
		relPaths = append(relPaths, relPath)
	}

	return relPaths, nil
}
