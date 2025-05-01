package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nikos/godora/internal/domain"
)

type GitAdapter struct {
	repoPath string
}

func NewGitAdapter(repoPath string) *GitAdapter {
	return &GitAdapter{
		repoPath: repoPath,
	}
}

func (g *GitAdapter) GetCommits(path string, startDate, endDate time.Time) ([]domain.Commit, error) {
	// Format dates for git command
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	// Execute git log command with the provided path
	cmd := exec.Command("git", "-C", path, "log", "--first-parent", "master", "--date=iso", "--since="+startDateStr, "--until="+endDateStr)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git log: %v", err)
	}

	return g.parseGitLog(bytes.NewReader(output))
}

func (g *GitAdapter) parseGitLog(reader io.Reader) ([]domain.Commit, error) {
	var commits []domain.Commit
	var currentCommit domain.Commit
	var commitMessage strings.Builder

	// Regular expressions for parsing
	commitRegex := regexp.MustCompile(`^commit ([0-9a-z]{40})`)
	authorRegex := regexp.MustCompile(`^Author:(.+)$`)
	dateRegex := regexp.MustCompile(`^Date:\s{0,10}([0-9]{4}-[0-9]{2}-[0-9]{2})`)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for commit line
		if matches := commitRegex.FindStringSubmatch(line); matches != nil {
			// If we have a previous commit, add it to the list
			if currentCommit.Hash != "" {
				// Process the accumulated message
				currentCommit.Hash = strings.TrimSpace(currentCommit.Hash)
				message := strings.TrimSpace(commitMessage.String())
				currentCommit.IsFeature = strings.Contains(message, "feat")
				currentCommit.IsFix = strings.Contains(message, "fix")
				currentCommit.IsDeploy = strings.Contains(message, "deploy")
				commits = append(commits, currentCommit)
				commitMessage.Reset()
			}

			// Start new commit
			currentCommit = domain.Commit{
				Hash:       strings.TrimSpace(matches[1]),
				BranchName: "main", // All commits are from master branch
			}
			continue
		}

		// Check for author line
		if matches := authorRegex.FindStringSubmatch(line); matches != nil {
			continue // We don't need author info
		}

		// Check for date line
		if matches := dateRegex.FindStringSubmatch(line); matches != nil {
			dateStr := strings.TrimSpace(matches[1])

			// split the dateStr by dash and create a new date
			dateParts := strings.Split(dateStr, "-")
			year, _ := strconv.Atoi(dateParts[0])
			month, _ := strconv.Atoi(dateParts[1])
			day, _ := strconv.Atoi(dateParts[2])
			date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			currentCommit.Date = date
			continue
		}

		// Accumulate message lines (skip empty lines)
		if strings.TrimSpace(line) != "" {
			if commitMessage.Len() > 0 {
				commitMessage.WriteString(" ")
			}
			commitMessage.WriteString(strings.TrimSpace(line))
		}
	}

	// Add the last commit
	if currentCommit.Hash != "" {
		message := strings.TrimSpace(commitMessage.String())
		currentCommit.IsFeature = strings.HasPrefix(message, "feat:")
		currentCommit.IsFix = strings.HasPrefix(message, "fix:")
		currentCommit.IsDeploy = strings.HasPrefix(message, "deploy:")
		commits = append(commits, currentCommit)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading git log: %v", err)
	}

	return commits, nil
}

func (g *GitAdapter) calculateMetrics(commits []domain.Commit) domain.DORAMetrics {
	if len(commits) == 0 {
		return domain.DORAMetrics{}
	}

	var (
		deployments      int
		failedDeploys    int
		featureLeadTimes []float64
	)

	// Calculate metrics
	for i := 0; i < len(commits); i++ {
		commit := commits[i]

		if commit.IsDeploy {
			deployments++
			// Check if the next commit is a fix
			if i+1 < len(commits) && commits[i+1].IsFix {
				failedDeploys++
			}
		}

		if commit.IsFeature {
			// Find the next deployment
			for j := i + 1; j < len(commits); j++ {
				if commits[j].IsDeploy {
					leadTime := commits[j].Date.Sub(commit.Date).Hours()
					featureLeadTimes = append(featureLeadTimes, leadTime)
					break
				}
			}
		}
	}

	// Calculate deployment frequency (deployments per day)
	days := commits[0].Date.Sub(commits[len(commits)-1].Date).Hours() / 24
	deploymentFrequency := float64(deployments) / days

	// Calculate change failure rate
	changeFailureRate := 0.0
	if deployments > 0 {
		changeFailureRate = float64(failedDeploys) / float64(deployments) * 100
	}

	// Calculate lead time for changes (average in hours)
	leadTimeForChanges := 0.0
	if len(featureLeadTimes) > 0 {
		sum := 0.0
		for _, lt := range featureLeadTimes {
			sum += lt
		}
		leadTimeForChanges = sum / float64(len(featureLeadTimes))
	}

	return domain.DORAMetrics{
		DeploymentFrequency: deploymentFrequency,
		LeadTimeForChanges:  leadTimeForChanges,
		ChangeFailureRate:   changeFailureRate,
	}
}
