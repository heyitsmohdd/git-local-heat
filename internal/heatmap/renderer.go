package heatmap

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/shariqattar/git-local-heat/pkg/models"
)

type Renderer struct {
	colors []string
}

func NewRenderer(colors []string) *Renderer {
	return &Renderer{colors: colors}
}

func (r *Renderer) Render(grid [][]models.HeatmapCell, commits models.CommitsByDate) string {
	var output strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))

	output.WriteString(headerStyle.Render("Git Contribution Heatmap") + "\n\n")

	totalCommits := GetTotalCommits(commits)
	maxCommits := GetMaxCommitCount(commits)
	statsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	output.WriteString(statsStyle.Render(fmt.Sprintf("Total commits: %d  |  Max in a day: %d", totalCommits, maxCommits)) + "\n\n")

	output.WriteString(r.renderMonthLabels(grid))
	output.WriteString("\n")

	dayLabels := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	for day := 0; day < 7; day++ {
		if day == 1 || day == 3 || day == 5 {
			labelStyle := lipgloss.NewStyle().
				Width(4).
				Foreground(lipgloss.Color("240"))
			output.WriteString(labelStyle.Render(dayLabels[day]))
		} else {
			output.WriteString("    ")
		}

		for week := 0; week < len(grid); week++ {
			if day < len(grid[week]) {
				cell := grid[week][day]
				output.WriteString(r.renderCell(cell))
			}
		}
		output.WriteString("\n")
	}

	output.WriteString("\n")
	output.WriteString(r.renderLegend())

	return output.String()
}

func (r *Renderer) renderCell(cell models.HeatmapCell) string {
	block := "■"
	color := r.hexToLipglossColor(r.colors[cell.ColorLevel])

	cellStyle := lipgloss.NewStyle().
		Foreground(color).
		Width(2)

	return cellStyle.Render(block + " ")
}

func (r *Renderer) renderMonthLabels(grid [][]models.HeatmapCell) string {
	var labels strings.Builder
	labels.WriteString("    ")

	currentMonth := -1
	monthStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Bold(true)

	for week := 0; week < len(grid); week++ {
		if len(grid[week]) > 0 {
			month := int(grid[week][0].Date.Month())

			if month != currentMonth {
				monthName := grid[week][0].Date.Format("Jan")
				labels.WriteString(monthStyle.Render(monthName))
				currentMonth = month
				labels.WriteString(" ")
			} else {
				labels.WriteString("  ")
			}
		}
	}

	return labels.String()
}

func (r *Renderer) renderLegend() string {
	var legend strings.Builder

	legendStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	legend.WriteString(legendStyle.Render("Less "))

	levels := []string{"■", "■", "■", "■", "■"}
	for i, block := range levels {
		color := r.hexToLipglossColor(r.colors[i])
		cellStyle := lipgloss.NewStyle().Foreground(color)
		legend.WriteString(cellStyle.Render(block + " "))
	}

	legend.WriteString(legendStyle.Render(" More"))

	return legend.String()
}

func (r *Renderer) hexToLipglossColor(hex string) lipgloss.Color {
	hex = strings.TrimPrefix(hex, "#")
	return lipgloss.Color("#" + hex)
}
