package api

import (
	"lib/internal/database"
	"net/http"
)

func HandleBookGet(db *[]database.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only GET Requests are allowed")
			return
		} else if r.PathValue("id") != "" {
			handleGetBook(db)(w, r)
			return
		}
		

		if len(r.URL.Query()) > 0 {
			handleSearchBook(db)(w, r)
			return

		}

		handleListBooks(db)(w, r)
	}
}
