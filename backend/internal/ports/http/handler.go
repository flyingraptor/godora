package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nikos/godora/internal/domain"
)

type MetricsHandler struct {
	metricsRepo domain.MetricsRepository
}

func NewMetricsHandler(metricsRepo domain.MetricsRepository) *MetricsHandler {
	return &MetricsHandler{
		metricsRepo: metricsRepo,
	}
}

type MetricsRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	RepoPath  string `json:"repoPath"`
}

func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	var req MetricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RepoPath == "" {
		http.Error(w, "Repository path is required", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}

	commits, err := h.metricsRepo.GetCommits(req.RepoPath, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to get commits", http.StatusInternalServerError)
		return
	}

	metrics := h.metricsRepo.CalculateMetrics(commits)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
