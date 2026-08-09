package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type RequestData struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Data    string `json:"data"`
	Version string `json:"version"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "observability-lab"
	}

	http.HandleFunc("/api/data", dataHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/metrics", metricsHandler)

	log.Printf("Service %s starting on port %s", serviceName, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Timestamp == "" {
		req.Timestamp = time.Now().Format(time.RFC3339)
	}

	log.Printf("Received - ID: %s, Message: %s, Time: %s", 
		req.ID, req.Message, req.Timestamp)

	resp := ResponseData{
		Status:  "success",
		Data:    req.Message,
		Version: "v1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

	log.Printf("Response sent in %v", time.Since(startTime))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total 42
`))
}
