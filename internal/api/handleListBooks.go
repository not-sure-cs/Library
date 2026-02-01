package api

import (
	"lib/internal/database"
	"net/http"
)

func handleListBooks(db *[]database.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		RespondWithJSON(w, http.StatusOK, db)
	}
}
