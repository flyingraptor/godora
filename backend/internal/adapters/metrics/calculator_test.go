package metrics

import (
	"testing"
	"time"

	"github.com/nikos/godora/internal/domain"
)

func TestCalculateMetrics(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		commits  []domain.Commit
		expected domain.DORAMetrics
	}{
		{
			name: "No deployments",
			commits: []domain.Commit{
				{
					Hash:      "1",
					Date:      now,
					IsFeature: true,
				},
			},
			expected: domain.DORAMetrics{
				DeploymentFrequency: 0,
				LeadTimeForChanges:  0,
				ChangeFailureRate:   0,
			},
		},
		{
			name: "Single deployment",
			commits: []domain.Commit{
				{
					Hash:      "1",
					Date:      now,
					IsFeature: true,
				},
				{
					Hash:     "2",
					Date:     now.Add(24 * time.Hour),
					IsDeploy: true,
				},
			},
			expected: domain.DORAMetrics{
				DeploymentFrequency: 1,
				LeadTimeForChanges:  0,
				ChangeFailureRate:   0,
			},
		},
		{
			name: "Deployment with fixes",
			commits: []domain.Commit{
				{
					Hash:      "1",
					Date:      now,
					IsFeature: true,
				},
				{
					Hash:     "2",
					Date:     now.Add(24 * time.Hour),
					IsDeploy: true,
				},
				{
					Hash:  "3",
					Date:  now.Add(25 * time.Hour),
					IsFix: true,
				},
				{
					Hash:  "4",
					Date:  now.Add(26 * time.Hour),
					IsFix: true,
				},
			},
			expected: domain.DORAMetrics{
				DeploymentFrequency: 1,
				LeadTimeForChanges:  0,
				ChangeFailureRate:   200,
			},
		},
		{
			name: "Multiple features and deployments",
			commits: []domain.Commit{
				{
					Hash:      "1",
					Date:      now,
					IsFeature: true,
				},
				{
					Hash:     "2",
					Date:     now.Add(24 * time.Hour),
					IsDeploy: true,
				},
				{
					Hash:      "3",
					Date:      now.Add(24 * time.Hour),
					IsFeature: true,
				},
				{
					Hash:     "4",
					Date:     now.Add(48 * time.Hour),
					IsDeploy: true,
				},
				{
					Hash:  "5",
					Date:  now.Add(49 * time.Hour),
					IsFix: true,
				},
			},
			expected: domain.DORAMetrics{
				DeploymentFrequency: 2,
				LeadTimeForChanges:  12, // (24h - 0h) / 2 features
				ChangeFailureRate:   50,
			},
		},
		{
			name: "Multiple features spread over time",
			commits: []domain.Commit{
				{
					Hash:      "1",
					Date:      now,
					IsFeature: true,
				},
				{
					Hash:      "2",
					Date:      now.Add(12 * time.Hour),
					IsFeature: true,
				},
				{
					Hash:      "3",
					Date:      now.Add(24 * time.Hour),
					IsFeature: true,
				},
				{
					Hash:     "4",
					Date:     now.Add(36 * time.Hour),
					IsDeploy: true,
				},
			},
			expected: domain.DORAMetrics{
				DeploymentFrequency: 1,
				LeadTimeForChanges:  8, // (24h - 0h) / 3 features
				ChangeFailureRate:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewMetricsCalculator()
			metrics := calculator.CalculateMetrics(tt.commits)

			if metrics.DeploymentFrequency != tt.expected.DeploymentFrequency {
				t.Errorf("DeploymentFrequency = %v, want %v", metrics.DeploymentFrequency, tt.expected.DeploymentFrequency)
			}
			if metrics.LeadTimeForChanges != tt.expected.LeadTimeForChanges {
				t.Errorf("LeadTimeForChanges = %v, want %v", metrics.LeadTimeForChanges, tt.expected.LeadTimeForChanges)
			}
			if metrics.ChangeFailureRate != tt.expected.ChangeFailureRate {
				t.Errorf("ChangeFailureRate = %v, want %v", metrics.ChangeFailureRate, tt.expected.ChangeFailureRate)
			}
		})
	}
}
