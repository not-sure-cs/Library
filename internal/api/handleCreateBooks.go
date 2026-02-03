package api

import (
	"encoding/json"
	"net/http"

	"github.com/knibirdgautam/library/internal/database"
)

func HandleCreateBooks(db *[]database.Book) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST requests allowed")
			return
		}

		var newBook database.Book

		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to decode POST request")
			return
		}

		newBook.ID = len(*db) + 1

		if !validateBook(newBook) {
			RespondWithError(w, http.StatusBadRequest, "Invalid Entry")
			return
		}

		*db = append(*db, newBook)

		RespondWithJSON(w, http.StatusOK, newBook)
	}
}
