package finder

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/koki-develop/go-fzf"
	"github.com/sourcegraph/conc/pool"
)

var ErrCancelled = errors.New("cancelled")

type entry struct {
	absPath      string
	relPath      string
	repoName     string
	modTime      time.Time
	displayLabel string
}

func Run(dirname string) (string, error) {
	repos, err := ghqRepos()
	if err != nil {
		return "", err
	}

	entries, err := collectFiles(repos, dirname)
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "", nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].modTime.After(entries[j].modTime)
	})

	for i := range entries {
		entries[i].displayLabel = fmt.Sprintf("%s  %s  [%s]",
			entries[i].modTime.Format("2006-01-02"),
			entries[i].relPath,
			entries[i].repoName,
		)
	}

	f, err := fzf.New(
		fzf.WithInputPosition(fzf.InputPositionTop),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create fuzzy finder: %w", err)
	}

	idxs, err := f.Find(
		entries,
		func(i int) string {
			return entries[i].displayLabel
		},
	)
	if err != nil {
		if errors.Is(err, fzf.ErrAbort) {
			return "", ErrCancelled
		}
		return "", fmt.Errorf("fuzzy finder error: %w", err)
	}

	selected := entries[idxs[0]]

	now := time.Now()
	if err := os.Chtimes(selected.absPath, now, now); err != nil {
		return "", fmt.Errorf("failed to update timestamp: %w", err)
	}

	return selected.absPath, nil
}

func ghqRepos() ([]string, error) {
	cmd := exec.Command("ghq", "list", "--full-path")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run 'ghq list': %w (is ghq installed?)", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var repos []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			repos = append(repos, line)
		}
	}

	if len(repos) == 0 {
		return nil, errors.New("no repositories found via ghq list")
	}

	return repos, nil
}

func collectFiles(repos []string, dirname string) ([]entry, error) {
	p := pool.NewWithResults[[]entry]()

	for _, repo := range repos {
		p.Go(func() []entry {
			return collectFilesFromRepo(repo, dirname)
		})
	}

	results := p.Wait()

	var entries []entry
	for _, result := range results {
		entries = append(entries, result...)
	}

	return entries, nil
}

func collectFilesFromRepo(repo, dirname string) []entry {
	hidenDir := filepath.Join(repo, dirname)
	info, err := os.Stat(hidenDir)
	if err != nil || !info.IsDir() {
		return nil
	}

	repoName := filepath.Base(repo)
	var entries []entry

	_ = filepath.WalkDir(hidenDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(hidenDir, path)
		entries = append(entries, entry{
			absPath:  path,
			relPath:  relPath,
			repoName: repoName,
			modTime:  info.ModTime(),
		})
		return nil
	})

	return entries
}
