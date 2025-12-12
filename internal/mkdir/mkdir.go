package mkdir

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var ErrNotInGitRepo = errors.New("not in a git repository")

// Run creates a date-based directory in the hiden directory of the current git repository.
// Returns the relative path from repository root.
func Run(dirname string) (string, error) {
	absPath, relPath, err := EnsureDir(dirname)
	if err != nil {
		return "", err
	}
	_ = absPath
	return relPath, nil
}

// EnsureDir creates a date-based directory in the hiden directory of the current git repository.
// Returns the absolute path and relative path from repository root.
func EnsureDir(dirname string) (absPath string, relPath string, err error) {
	// Get git repository root
	repoRoot, err := getGitRepoRoot()
	if err != nil {
		return "", "", err
	}

	// Get current date in YYYY-MM-DD format
	today := time.Now().Format("2006-01-02")

	// Construct the directory path
	absPath = filepath.Join(repoRoot, dirname, today)

	// Create the directory (including parent directories if needed)
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Return relative path from repository root
	relPath = filepath.Join(dirname, today)

	return absPath, relPath, nil
}

// getGitRepoRoot returns the root directory of the git repository
func getGitRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		// Check if it's because we're not in a git repo
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 128 {
			return "", ErrNotInGitRepo
		}
		return "", fmt.Errorf("failed to get git repository root: %w", err)
	}

	repoRoot := strings.TrimSpace(string(output))
	return repoRoot, nil
}
