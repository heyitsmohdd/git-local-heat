package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/shariqattar/git-local-heat/pkg/models"
)

type Scanner struct {
	maxWorkers int
}

func NewScanner(maxWorkers int) *Scanner {
	if maxWorkers <= 0 {
		maxWorkers = 50
	}
	return &Scanner{maxWorkers: maxWorkers}
}

func (s *Scanner) ScanDirectory(root string) ([]models.Repository, error) {
	repoPaths := make(chan string, 100)
	results := make([]models.Repository, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < s.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range repoPaths {
				repo := s.validateRepository(path)
				mu.Lock()
				results = append(results, repo)
				mu.Unlock()
			}
		}()
	}

	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Warning: cannot access %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() && d.Name() == ".git" {
			repoPaths <- filepath.Dir(path)
			return fs.SkipDir
		}

		if d.IsDir() && (d.Name() == "node_modules" || d.Name() == "vendor" || d.Name() == ".cache") {
			return fs.SkipDir
		}

		return nil
	})

	close(repoPaths)
	wg.Wait()

	if walkErr != nil && len(results) == 0 {
		return nil, fmt.Errorf("failed to walk directory: %w", walkErr)
	}

	return results, nil
}

func (s *Scanner) validateRepository(path string) models.Repository {
	gitDir := filepath.Join(path, ".git")

	info, err := filepath.Glob(gitDir)
	if err != nil || len(info) == 0 {
		return models.Repository{
			Path:    path,
			IsValid: false,
			Error:   fmt.Errorf("invalid .git directory"),
		}
	}

	return models.Repository{
		Path:    path,
		IsValid: true,
		Error:   nil,
	}
}
