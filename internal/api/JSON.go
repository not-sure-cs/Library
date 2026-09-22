package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, msg string) {
	if statusCode >= http.StatusInternalServerError {
		log.Printf("CRITICAL ERROR %d, %s", statusCode, msg)
	}
	RespondWithJSON(w, statusCode, map[string]string{"error": msg})
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(payload)

	if err != nil {
		log.Printf("Failed to encode JSON response: %v", payload)
		// Note: We can't change the status code now because headers are already sent.
		return
	}
}
