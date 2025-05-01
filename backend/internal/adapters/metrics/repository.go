package metrics

import (
	"time"

	"github.com/nikos/godora/internal/adapters/git"
	"github.com/nikos/godora/internal/domain"
)

type MetricsRepository struct {
	gitAdapter *git.GitAdapter
	calculator *MetricsCalculator
}

func NewMetricsRepository(gitAdapter *git.GitAdapter, calculator *MetricsCalculator) *MetricsRepository {
	return &MetricsRepository{
		gitAdapter: gitAdapter,
		calculator: calculator,
	}
}

func (r *MetricsRepository) GetCommits(path string, startDate, endDate time.Time) ([]domain.Commit, error) {
	return r.gitAdapter.GetCommits(path, startDate, endDate)
}

func (r *MetricsRepository) CalculateMetrics(commits []domain.Commit) domain.DORAMetrics {
	return r.calculator.CalculateMetrics(commits)
}
