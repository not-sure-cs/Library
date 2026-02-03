package api

import (
	"net/http"

	"github.com/knibirdgautam/library/internal/database"
)

func handleListBooks(db *[]database.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		RespondWithJSON(w, http.StatusOK, db)
	}
}
