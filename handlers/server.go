package handlers

import (
	"encoding/json"
	"net/http"

	"webserver/kafka"
	"webserver/redis"
)

// ProduceHandler handles POST /produce to send messages to Kafka.
func ProduceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data redis.Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if data.ID == "" || data.Value == "" {
		http.Error(w, "Missing id or value", http.StatusBadRequest)
		return
	}

	if err := kafka.ProduceMessage(data); err != nil {
		http.Error(w, "Failed to produce message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Message produced"})
}