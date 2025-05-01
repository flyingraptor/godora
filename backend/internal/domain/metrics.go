package domain

import "time"

// Commit represents a git commit with its metadata
type Commit struct {
	Hash       string
	Author     string
	Date       time.Time
	Message    string
	BranchName string
	IsDeploy   bool
	IsFeature  bool
	IsFix      bool
}

// DORAMetrics represents the calculated DORA metrics
type DORAMetrics struct {
	DeploymentFrequency float64
	LeadTimeForChanges  float64
	ChangeFailureRate   float64
}

// MetricsRepository defines the interface for fetching metrics
type MetricsRepository interface {
	GetCommits(path string, startDate, endDate time.Time) ([]Commit, error)
	CalculateMetrics(commits []Commit) DORAMetrics
}
