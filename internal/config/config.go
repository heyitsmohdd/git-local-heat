package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shariqattar/git-local-heat/pkg/models"
	"gopkg.in/yaml.v3"
)

var DefaultColors = []string{
	"#161b22",
	"#0e4429",
	"#006d32",
	"#26a641",
	"#39d353",
}

func Load(configPath, email, path string, months int) (*models.Config, error) {
	config := &models.Config{
		Email:       email,
		Path:        path,
		Months:      months,
		ColorScheme: DefaultColors,
	}

	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configPath = filepath.Join(homeDir, ".git-local-heat.yaml")
		}
	}

	if configPath != "" && fileExists(configPath) {
		fileConfig, err := loadFromFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}

		if config.Email == "" {
			config.Email = fileConfig.Email
		}
		if config.Path == "" || config.Path == "." {
			config.Path = fileConfig.Path
		}
		if config.Months == 0 {
			config.Months = fileConfig.Months
		}
		if len(fileConfig.ColorScheme) > 0 {
			config.ColorScheme = fileConfig.ColorScheme
		}
	}

	if config.Path == "" || config.Path == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		config.Path = cwd
	}

	if config.Months == 0 {
		config.Months = 12
	}

	if config.Email == "" {
		config.Email = getGitConfigEmail()
	}

	return config, nil
}

func loadFromFile(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func getGitConfigEmail() string {
	cmd := exec.Command("git", "config", "--global", "user.email")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
