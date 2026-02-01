package api

import (
	"net/http"
	"time"
)

func HandleStatus(start time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(start)
		RespondWithJSON(w, http.StatusOK, map[string]time.Duration{"Uptime": uptime})
	}
}
