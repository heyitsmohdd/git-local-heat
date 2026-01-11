package heatmap

import (
	"time"

	"github.com/shariqattar/git-local-heat/pkg/models"
)

type Generator struct {
	months int
}

func NewGenerator(months int) *Generator {
	return &Generator{months: months}
}

func (g *Generator) GenerateGrid(commits models.CommitsByDate) [][]models.HeatmapCell {
	endDate := time.Now()
	startDate := endDate.AddDate(0, -g.months, 0)

	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	days := int(endDate.Sub(startDate).Hours() / 24)
	weeks := (days / 7) + 1

	grid := make([][]models.HeatmapCell, weeks)
	for i := range grid {
		grid[i] = make([]models.HeatmapCell, 7)
	}

	currentDate := startDate
	for week := 0; week < weeks; week++ {
		for day := 0; day < 7; day++ {
			if currentDate.After(endDate) {
				break
			}

			dateKey := currentDate.Format("2006-01-02")
			count := commits[dateKey]

			grid[week][day] = models.HeatmapCell{
				Date:       currentDate,
				Count:      count,
				ColorLevel: calculateColorLevel(count),
			}

			currentDate = currentDate.AddDate(0, 0, 1)
		}
	}

	return grid
}

func calculateColorLevel(count int) int {
	switch {
	case count == 0:
		return 0
	case count <= 2:
		return 1
	case count <= 5:
		return 2
	case count <= 10:
		return 3
	default:
		return 4
	}
}

func GetMaxCommitCount(commits models.CommitsByDate) int {
	max := 0
	for _, count := range commits {
		if count > max {
			max = count
		}
	}
	return max
}

func GetTotalCommits(commits models.CommitsByDate) int {
	total := 0
	for _, count := range commits {
		total += count
	}
	return total
}
