package api

import (
	"net/http"
	"strconv"

	"github.com/knibirdgautam/library/internal/database"
)

func handleGetBook(db *[]database.Book) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.PathValue("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			RespondWithError(w, http.StatusUnprocessableEntity, "ID is invalid")
			return
		}

		book, err1 := findBook(db, id)
		if err1 != nil {
			RespondWithError(w, http.StatusNotFound, err1.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, *book)
	}
}
