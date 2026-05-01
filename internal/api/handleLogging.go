package api

import (
	"encoding/json"
	"net/http"

	"github.com/knibirdgautam/library/internal/database"
)

func HandleLogging(queries database.DBQueries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST requests allowed")
			return
		}

		type parameters struct {
			Email    string `json:"email"`
			password string `json:"password"`
		}

		var params parameters

		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Couldn't Parse POST request")
		}

		

		RespondWithJSON(w,400,pass)

	}

}
