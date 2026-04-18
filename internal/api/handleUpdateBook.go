package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

func HandleUpdateBooks(queries database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT Requests allowed")
			return
		}

		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)

		if err != nil {
			RespondWithError(w, http.StatusUnprocessableEntity, "Couldn't Parse ID")
			return
		}

		decoder := json.NewDecoder(r.Body)
		params := database.Parameters{}
		err = decoder.Decode(&params)

		book, err := queries.GetBook(r.Context(), id)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		

		RespondWithJSON(w, http.StatusOK, book)

	}
}
