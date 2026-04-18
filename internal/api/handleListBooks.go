package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

func HandleListOfAuthorBooks(queries *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only GET requests allowed")
			return
		}

		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)

		if err != nil {
			RespondWithError(w, http.StatusUnprocessableEntity, "Couldn't Parse ID")
			return
		}

		
		books,err := queries.GetAuthorBooks(r.Context(),id)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Unable to Find author")
			return 
		}

		RespondWithJSON(w, http.StatusOK, books)
	}
}
