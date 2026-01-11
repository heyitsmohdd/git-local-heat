package parser

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/shariqattar/git-local-heat/pkg/models"
)

type Parser struct {
	email  string
	months int
}

func NewParser(email string, months int) *Parser {
	return &Parser{
		email:  email,
		months: months,
	}
}

func (p *Parser) ParseRepositories(repos []models.Repository) (models.CommitsByDate, error) {
	commits := make(models.CommitsByDate)
	var mu sync.Mutex
	var wg sync.WaitGroup

	cutoffDate := time.Now().AddDate(0, -p.months, 0)

	for _, repo := range repos {
		if !repo.IsValid {
			continue
		}

		wg.Add(1)
		go func(repoPath string) {
			defer wg.Done()

			repoCommits, err := p.parseRepository(repoPath, cutoffDate)
			if err != nil {
				fmt.Printf("Warning: failed to parse %s: %v\n", repoPath, err)
				return
			}

			mu.Lock()
			for date, count := range repoCommits {
				commits[date] += count
			}
			mu.Unlock()
		}(repo.Path)
	}

	wg.Wait()
	return commits, nil
}

func (p *Parser) parseRepository(repoPath string, cutoffDate time.Time) (models.CommitsByDate, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	commits := make(models.CommitsByDate)

	err = commitIter.ForEach(func(c *object.Commit) error {
		if c.Author.When.Before(cutoffDate) {
			return nil
		}

		if p.email != "" && c.Author.Email != p.email {
			return nil
		}

		dateKey := c.Author.When.Format("2006-01-02")
		commits[dateKey]++

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate commits: %w", err)
	}

	return commits, nil
}
