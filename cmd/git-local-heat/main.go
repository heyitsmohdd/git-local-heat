package main

import (
	"fmt"
	"os"

	"github.com/shariqattar/git-local-heat/internal/config"
	"github.com/shariqattar/git-local-heat/internal/heatmap"
	"github.com/shariqattar/git-local-heat/internal/parser"
	"github.com/shariqattar/git-local-heat/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	email      string
	path       string
	months     int
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "git-local-heat",
	Short: "Visualize git commit history as a terminal heatmap",
	Long: `git-local-heat scans directories for git repositories and displays
a GitHub-style contribution heatmap in your terminal.

It recursively finds all .git repositories, extracts commit history,
and generates a beautiful visual representation of your coding activity.`,
	RunE: runHeatmap,
}

func init() {
	rootCmd.Flags().StringVarP(&email, "email", "e", "", "Filter commits by author email (default: git config user.email)")
	rootCmd.Flags().StringVarP(&path, "path", "p", ".", "Root directory to scan (default: current directory)")
	rootCmd.Flags().IntVarP(&months, "months", "m", 12, "Number of months to display (6 or 12)")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Config file path (default: ~/.git-local-heat.yaml)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runHeatmap(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(configFile, email, path, months)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Display configuration
	fmt.Printf("🔍 Scanning: %s\n", cfg.Path)
	if cfg.Email != "" {
		fmt.Printf("📧 Email filter: %s\n", cfg.Email)
	} else {
		fmt.Println("📧 Email filter: all authors")
	}
	fmt.Printf("📅 Timeframe: %d months\n\n", cfg.Months)

	// Step 1: Scan for repositories
	fmt.Println("⏳ Scanning for git repositories...")
	s := scanner.NewScanner(50)
	repos, err := s.ScanDirectory(cfg.Path)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	validRepos := 0
	for _, repo := range repos {
		if repo.IsValid {
			validRepos++
		}
	}

	fmt.Printf("✅ Found %d repositories\n\n", validRepos)

	if validRepos == 0 {
		fmt.Println("No git repositories found. Try a different path.")
		return nil
	}

	// Step 2: Parse commit history
	fmt.Println("⏳ Parsing commit history...")
	p := parser.NewParser(cfg.Email, cfg.Months)
	commits, err := p.ParseRepositories(repos)
	if err != nil {
		return fmt.Errorf("failed to parse repositories: %w", err)
	}

	fmt.Printf("✅ Processed commits\n\n")

	if len(commits) == 0 {
		fmt.Println("No commits found matching the criteria.")
		return nil
	}

	// Step 3: Generate heatmap grid
	g := heatmap.NewGenerator(cfg.Months)
	grid := g.GenerateGrid(commits)

	// Step 4: Render heatmap
	r := heatmap.NewRenderer(cfg.ColorScheme)
	output := r.Render(grid, commits)

	fmt.Println(output)

	return nil
}
