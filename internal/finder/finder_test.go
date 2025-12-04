package finder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectFilesFromRepo_WithSymlink(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "test-repo")
	realHidenDir := filepath.Join(repoDir, ".hiden_real")
	symlinkHidenDir := filepath.Join(repoDir, ".hiden")

	// Create repo directory
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	// Create real hiden directory with a test file
	if err := os.MkdirAll(realHidenDir, 0755); err != nil {
		t.Fatalf("Failed to create real hiden dir: %v", err)
	}
	testFile := filepath.Join(realHidenDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create subdirectory with another file
	subDir := filepath.Join(realHidenDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	nestedFile := filepath.Join(subDir, "nested.txt")
	if err := os.WriteFile(nestedFile, []byte("nested content"), 0644); err != nil {
		t.Fatalf("Failed to create nested file: %v", err)
	}

	// Create symlink
	if err := os.Symlink(".hiden_real", symlinkHidenDir); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	// Debug: Check symlink
	t.Logf("Repo dir: %s", repoDir)
	t.Logf("Symlink path: %s", symlinkHidenDir)
	info, err := os.Lstat(symlinkHidenDir)
	if err != nil {
		t.Logf("Lstat error: %v", err)
	} else {
		t.Logf("Is symlink: %v", info.Mode()&os.ModeSymlink != 0)
	}

	// Test collectFilesFromRepo with symlinked hiden directory
	entries := collectFilesFromRepo(repoDir, ".hiden")

	// Should find 2 files (test.txt and nested.txt)
	t.Logf("Found %d entries", len(entries))
	for _, e := range entries {
		t.Logf("Found: %s (rel: %s)", e.absPath, e.relPath)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 files, got %d", len(entries))
	}

	// Verify that files are found
	foundFiles := make(map[string]bool)
	for _, e := range entries {
		foundFiles[e.relPath] = true
	}

	if !foundFiles["test.txt"] {
		t.Error("Expected to find test.txt")
	}
	if !foundFiles["subdir/nested.txt"] && !foundFiles["subdir\\nested.txt"] {
		t.Error("Expected to find subdir/nested.txt")
	}
}

func TestCollectFilesFromRepo_WithoutSymlink(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "test-repo")
	hidenDir := filepath.Join(repoDir, ".hiden")

	// Create repo directory
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	// Create regular hiden directory with a test file
	if err := os.MkdirAll(hidenDir, 0755); err != nil {
		t.Fatalf("Failed to create hiden dir: %v", err)
	}
	testFile := filepath.Join(hidenDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test collectFilesFromRepo with regular hiden directory
	entries := collectFilesFromRepo(repoDir, ".hiden")

	// Should find 1 file
	if len(entries) != 1 {
		t.Errorf("Expected 1 file, got %d", len(entries))
	}

	if len(entries) > 0 && entries[0].relPath != "test.txt" {
		t.Errorf("Expected relPath to be 'test.txt', got '%s'", entries[0].relPath)
	}
}
