package api

import (
	"net/http"

	"github.com/knibirdgautam/library/internal/database"
)

func HandleListOfAuthorBooks(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only GET requests allowed")
			return
		}

		author := r.PathValue("author")
		books,err := queries.GetAuthorBook(r.Context(),author)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Unable to Find author")
			return 
		}

		RespondWithJSON(w, http.StatusOK, books)
	}
}
