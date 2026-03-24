package api

import (
	"net/http"

	"github.com/knibirdgautam/library/internal/database"
)

func HandleGetBooks(queries *database.Queries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only GET requests allowed")
			return
		}

		name := r.PathValue("name")

		book,err := queries.GetBook(r.Context(), name)
		if err!= nil{
			RespondWithError(w, http.StatusBadRequest, "Unable to Find Name")
		}

		RespondWithJSON(w, http.StatusOK,book)
	}
}
