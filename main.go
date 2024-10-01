package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	healthStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_app_health_status",
		Help: "The health status of the Go app (1 = healthy, 0 = unhealthy)",
	})
)

func init() {
	prometheus.MustRegister(healthStatus)
}

func main() {
	// Set health to healthy on startup
	healthStatus.Set(1)

	// Handle the /health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "healthy")
		healthStatus.Set(1) // App is healthy when this endpoint is hit
	})

	// Handle the /ping endpoint
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong")
	})

	// Handle the /unhealthy endpoint
	// This endpoint will set the health status to unhealthy
	http.HandleFunc("/unhealthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "unhealthy")
		healthStatus.Set(0) // App is unhealthy when this endpoint is hit
	})

	// Expose the Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start the server on port 8080
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
