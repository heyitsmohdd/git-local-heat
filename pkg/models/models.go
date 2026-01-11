package models

import "time"

type Commit struct {
	Hash    string
	Author  string
	Email   string
	Date    time.Time
	Message string
}

type Config struct {
	Email       string
	Path        string
	Months      int
	ColorScheme []string
}

type HeatmapCell struct {
	Date       time.Time
	Count      int
	ColorLevel int
}

type CommitsByDate map[string]int

type Repository struct {
	Path    string
	IsValid bool
	Error   error
}
