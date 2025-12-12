package mv

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/qawatake/hiden/internal/mkdir"
)

func TestRun_Success(t *testing.T) {
	// Create temporary directory for git repo
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Create a test file to move
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to the git repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Run the mv command
	result, err := Run(".hiden", testFile)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify the result path format
	today := time.Now().Format("2006-01-02")
	expectedRelPath := filepath.Join(".hiden", today, "test.txt")
	if result != expectedRelPath {
		t.Errorf("Expected result %q, got %q", expectedRelPath, result)
	}

	// Verify the file was moved
	expectedAbsPath := filepath.Join(tmpDir, ".hiden", today, "test.txt")
	if _, err := os.Stat(expectedAbsPath); os.IsNotExist(err) {
		t.Errorf("File was not moved to expected location: %s", expectedAbsPath)
	}

	// Verify the original file no longer exists
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Original file still exists after move")
	}

	// Verify file content
	content, err := os.ReadFile(expectedAbsPath)
	if err != nil {
		t.Fatalf("Failed to read moved file: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("File content mismatch: expected %q, got %q", "test content", string(content))
	}
}

func TestRun_NotInGitRepo(t *testing.T) {
	// Create temporary directory (not a git repo)
	tmpDir := t.TempDir()

	// Create a test file to move
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to the non-git directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Run the mv command
	_, err = Run(".hiden", testFile)
	if err == nil {
		t.Fatal("Expected error when not in git repo")
	}
	if err != mkdir.ErrNotInGitRepo {
		t.Errorf("Expected mkdir.ErrNotInGitRepo, got: %v", err)
	}
}

func TestRun_FileNotExist(t *testing.T) {
	// Create temporary directory for git repo
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Change to the git repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Try to move a non-existent file
	_, err = Run(".hiden", filepath.Join(tmpDir, "nonexistent.txt"))
	if err == nil {
		t.Fatal("Expected error when file does not exist")
	}
	if !strings.Contains(err.Error(), "failed to move file") {
		t.Errorf("Expected 'failed to move file' error, got: %v", err)
	}
}

func TestRun_CreatesDirectory(t *testing.T) {
	// Create temporary directory for git repo
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Create a test file to move
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify the hiden directory doesn't exist yet
	today := time.Now().Format("2006-01-02")
	hidenDir := filepath.Join(tmpDir, ".hiden", today)
	if _, err := os.Stat(hidenDir); !os.IsNotExist(err) {
		t.Fatal("Hiden directory should not exist before mv")
	}

	// Change to the git repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Run the mv command
	_, err = Run(".hiden", testFile)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify the directory was created
	if _, err := os.Stat(hidenDir); os.IsNotExist(err) {
		t.Error("Hiden directory was not created")
	}
}
