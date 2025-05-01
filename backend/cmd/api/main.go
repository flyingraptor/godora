package main

import (
	"log"
	"net/http"

	"github.com/nikos/godora/internal/adapters/git"
	"github.com/nikos/godora/internal/adapters/metrics"
	httphandler "github.com/nikos/godora/internal/ports/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize adapters
	gitAdapter := git.NewGitAdapter(".") // Assuming the git repo is in the current directory
	calcAdapter := metrics.NewMetricsCalculator()
	metricsRepo := metrics.NewMetricsRepository(gitAdapter, calcAdapter)

	// Initialize handler
	handler := httphandler.NewMetricsHandler(metricsRepo)

	// Setup routes with CORS middleware
	http.Handle("/api/metrics", corsMiddleware(http.HandlerFunc(handler.GetMetrics)))

	// Serve frontend static files
	fs := http.FileServer(http.Dir("../../frontend/build"))
	http.Handle("/", fs)

	// Start server
	log.Println("Server starting on :8099")
	if err := http.ListenAndServe(":8099", nil); err != nil {
		log.Fatal(err)
	}
}
