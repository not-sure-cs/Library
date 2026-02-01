package api

import (
	"encoding/json"
	"lib/internal/database"
	"net/http"
	"strconv"
)

func HandleUpdateBooks(db *[]database.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT Requests allowed")
			return
		}

		idStr := r.PathValue("id")

		var newBook database.Book

		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to decode PUT request")
			return
		}

		if !validateBook(newBook) {
			RespondWithError(w, http.StatusBadRequest, "Invalid Entry")
			return
		}

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

		book.Title = newBook.Title
		book.Author = newBook.Author

		RespondWithJSON(w, http.StatusOK, *book)

	}
}
