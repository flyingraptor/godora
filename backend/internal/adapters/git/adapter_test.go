package git

import (
	"strings"
	"testing"
	"time"

	"github.com/nikos/godora/internal/domain"
)

func TestParseGitLog(t *testing.T) {
	// Sample git log output matching git log --first-parent master format
	sampleLog := `commit 7a8b9c0d1e2f3g4h5i6j7k8l9m0n1o2p3q4r5s6t7u8v9w0x1y2z3
Author: John Doe <john.doe@example.com>
Date:   Thu Dec 28 13:43:41 2023 +0000

    feat: implement user authentication

commit 1d2e3f4g5h6i7j8k9l0m1n2o3p4q5r6s7t8u9v0w1x2y3z4a5b6c7d8e9f0
Author: Jane Smith <jane.smith@example.com>
Date:   Thu Dec 28 12:30:15 2023 +0000

    fix: resolve login timeout issue

commit a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0
Author: John Doe <john.doe@example.com>
Date:   Thu Dec 28 10:15:22 2023 +0000

    deploy: release v1.0.0

commit 4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2g3h4
Author: Jane Smith <jane.smith@example.com>
Date:   Thu Dec 27 16:45:33 2023 +0000

    feat: add user profile management

commit 8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9
Author: John Doe <john.doe@example.com>
Date:   Thu Dec 27 14:20:10 2023 +0000

    fix: correct profile image upload

commit 2k3l4m5n6o7p8q9r0s1t2u3v4w5x6y7z8a9b0c1d2e3f4g5h6i7j8k9l0m1n2o3
Author: Jane Smith <jane.smith@example.com>
Date:   Thu Dec 27 11:05:45 2023 +0000

    feat: initial database setup

commit 6n7o8p9q0r1s2t3u4v5w6x7y8z9a0b1c2d3e4f5g6h7i8j9k0l1m2n3o4p5q6r7
Author: John Doe <john.doe@example.com>
Date:   Thu Dec 27 09:30:12 2023 +0000

    deploy: setup project structure`

	// Create test cases
	tests := []struct {
		name     string
		log      string
		expected []domain.Commit
	}{
		{
			name: "Parse sample git log",
			log:  sampleLog,
			expected: []domain.Commit{
				{
					Hash:       "7a8b9c0d1e2f3g4h5i6j7k8l9m0n1o2p3q4r5s6t7u8v9w0x1y2z3",
					BranchName: "main",
					IsDeploy:   false,
					IsFeature:  true,
					IsFix:      false,
					Date:       time.Date(2023, 12, 28, 13, 43, 41, 0, time.UTC),
				},
				{
					Hash:       "1d2e3f4g5h6i7j8k9l0m1n2o3p4q5r6s7t8u9v0w1x2y3z4a5b6c7d8e9f0",
					BranchName: "main",
					IsDeploy:   false,
					IsFeature:  false,
					IsFix:      true,
					Date:       time.Date(2023, 12, 28, 12, 30, 15, 0, time.UTC),
				},
				{
					Hash:       "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0",
					BranchName: "main",
					IsDeploy:   true,
					IsFeature:  false,
					IsFix:      false,
					Date:       time.Date(2023, 12, 28, 10, 15, 22, 0, time.UTC),
				},
				{
					Hash:       "4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2g3h4",
					BranchName: "main",
					IsDeploy:   false,
					IsFeature:  true,
					IsFix:      false,
					Date:       time.Date(2023, 12, 27, 16, 45, 33, 0, time.UTC),
				},
				{
					Hash:       "8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9",
					BranchName: "main",
					IsDeploy:   false,
					IsFeature:  false,
					IsFix:      true,
					Date:       time.Date(2023, 12, 27, 14, 20, 10, 0, time.UTC),
				},
				{
					Hash:       "2k3l4m5n6o7p8q9r0s1t2u3v4w5x6y7z8a9b0c1d2e3f4g5h6i7j8k9l0m1n2o3",
					BranchName: "main",
					IsDeploy:   false,
					IsFeature:  true,
					IsFix:      false,
					Date:       time.Date(2023, 12, 27, 11, 5, 45, 0, time.UTC),
				},
				{
					Hash:       "6n7o8p9q0r1s2t3u4v5w6x7y8z9a0b1c2d3e4f5g6h7i8j9k0l1m2n3o4p5q6r7",
					BranchName: "main",
					IsDeploy:   false,
					IsFeature:  false,
					IsFix:      false,
					Date:       time.Date(2023, 12, 27, 9, 30, 12, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewGitAdapter(".")
			commits, err := adapter.parseGitLog(strings.NewReader(tt.log))
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(commits) != len(tt.expected) {
				t.Fatalf("Expected %d commits, got %d", len(tt.expected), len(commits))
			}

			for i, commit := range commits {
				expected := tt.expected[i]
				if commit.Hash != expected.Hash {
					t.Errorf("Commit %d: Hash = %v, want %v", i, commit.Hash, expected.Hash)
				}
				if commit.BranchName != expected.BranchName {
					t.Errorf("Commit %d: BranchName = %v, want %v", i, commit.BranchName, expected.BranchName)
				}
				if commit.IsDeploy != expected.IsDeploy {
					t.Errorf("Commit %d: IsDeploy = %v, want %v", i, commit.IsDeploy, expected.IsDeploy)
				}
				if commit.IsFeature != expected.IsFeature {
					t.Errorf("Commit %d: IsFeature = %v, want %v", i, commit.IsFeature, expected.IsFeature)
				}
				if commit.IsFix != expected.IsFix {
					t.Errorf("Commit %d: IsFix = %v, want %v", i, commit.IsFix, expected.IsFix)
				}
				if !commit.Date.Equal(expected.Date) {
					t.Errorf("Commit %d: Date = %v, want %v", i, commit.Date, expected.Date)
				}
			}
		})
	}
}

func TestCalculateMetricsFromGitLog(t *testing.T) {
	// Use the same sample log as above
	sampleLog := `commit 7a8b9c0d1e2f3g4h5i6j7k8l9m0n1o2p3q4r5s6t7u8v9w0x1y2z3
Author: John Doe <john.doe@example.com>
Date:   Thu Dec 28 13:43:41 2023 +0000

    feat: implement user authentication

commit 1d2e3f4g5h6i7j8k9l0m1n2o3p4q5r6s7t8u9v0w1x2y3z4a5b6c7d8e9f0
Author: Jane Smith <jane.smith@example.com>
Date:   Thu Dec 28 12:30:15 2023 +0000

    fix: resolve login timeout issue

commit a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0
Author: John Doe <john.doe@example.com>
Date:   Thu Dec 28 10:15:22 2023 +0000

    deploy: release v1.0.0

commit 4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2g3h4
Author: Jane Smith <jane.smith@example.com>
Date:   Thu Dec 27 16:45:33 2023 +0000

    feat: add user profile management

commit 8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z7a8b9c0d1e2f3g4h5i6j7k8l9
Author: John Doe <john.doe@example.com>
Date:   Thu Dec 27 14:20:10 2023 +0000

    fix: correct profile image upload

commit 2k3l4m5n6o7p8q9r0s1t2u3v4w5x6y7z8a9b0c1d2e3f4g5h6i7j8k9l0m1n2o3
Author: Jane Smith <jane.smith@example.com>
Date:   Thu Dec 27 11:05:45 2023 +0000

    feat: initial database setup

commit 6n7o8p9q0r1s2t3u4v5w6x7y8z9a0b1c2d3e4f5g6h7i8j9k0l1m2n3o4p5q6r7
Author: John Doe <john.doe@example.com>
Date:   Thu Dec 27 09:30:12 2023 +0000

    deploy: setup project structure`

	adapter := NewGitAdapter(".")
	commits, err := adapter.parseGitLog(strings.NewReader(sampleLog))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Calculate metrics
	metrics := adapter.calculateMetrics(commits)

	// Verify metrics
	expected := domain.DORAMetrics{
		DeploymentFrequency: 0.5, // 1 deployment in 2 days
		LeadTimeForChanges:  3.5, // Average time from feature to main in hours
		ChangeFailureRate:   0.0, // No failed deployments to main
	}

	epsilon := 0.01
	if !floatEquals(metrics.DeploymentFrequency, expected.DeploymentFrequency, epsilon) {
		t.Errorf("DeploymentFrequency = %v, want %v", metrics.DeploymentFrequency, expected.DeploymentFrequency)
	}
	if !floatEquals(metrics.LeadTimeForChanges, expected.LeadTimeForChanges, epsilon) {
		t.Errorf("LeadTimeForChanges = %v, want %v", metrics.LeadTimeForChanges, expected.LeadTimeForChanges)
	}
	if !floatEquals(metrics.ChangeFailureRate, expected.ChangeFailureRate, epsilon) {
		t.Errorf("ChangeFailureRate = %v, want %v", metrics.ChangeFailureRate, expected.ChangeFailureRate)
	}
}

func floatEquals(a, b, epsilon float64) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}
