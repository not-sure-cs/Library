package api

import (
	"net/http"
	"strconv"

	"github.com/knibirdgautam/library/internal/database"
)

func HandleDeleteBook(db *[]database.Book) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE Requests are allowed")
			return
		}

		idStr := r.PathValue("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			RespondWithError(w, http.StatusUnprocessableEntity, "ID is invalid")
			return
		}

		err1 := delBook(db, id)
		if err1 != nil {
			RespondWithError(w, http.StatusNotFound, err1.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Book deleted"})
	}
}
