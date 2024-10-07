package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	// Define a CounterVec with labels for "url" and "success"
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_app_http_requests_total",
			Help: "Total number of HTTP requests received, labeled by URL and success status",
		},
		[]string{"url", "success"}, // Define labels
	)
)

func init() {
	prometheus.MustRegister(requestCounter)
}

func main() {
	// Initialize the Zap logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // Flushes buffer, if any

	// Handle the /health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "healthy")

		requestCounter.WithLabelValues("/health", "true").Inc()
	})

	// Handle the /ping endpoint
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")

		requestCounter.WithLabelValues("/ping", "true").Inc()
	})

	// Handle the /error endpoint to trigger an error log
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error occurred")

		requestCounter.WithLabelValues("/error", "false").Inc()

		// Log the error in JSON format with Zap
		logger.Error("An error occurred",
			zap.String("endpoint", "/error"),
			zap.String("severity", "error"),
			zap.String("message", "Triggered an error log"),
		)
	})

	// Expose the Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start the server on port 8080
	logger.Info("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
