package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

func HandleDeleteBook(queries database.DBQueries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE Requests are allowed")
			return
		}

		idStr := r.PathValue("id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			RespondWithError(w, http.StatusUnprocessableEntity, "Couldn't Parse ID")
			return
		}

		err = queries.UnlinkBook(r.Context(), id)

		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Couldn't Unlink book")
			return
		}

		err = queries.DeleteBook(r.Context(), id)

		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Couldn't Delete book")
			return
		}

		RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Book deleted"})
	}
}
