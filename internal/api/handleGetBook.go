package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
	"github.com/knibirdgautam/library/internal/storage"
)

func HandleGetBooks(queries database.DBQueries, store storage.R2Store, secret storage.Secret) http.HandlerFunc {

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

		book, err := queries.GetBook(r.Context(), id)
		if err != nil {
			RespondWithError(w, http.StatusNotFound, "Book not found")
			return
		}

		URL, err := store.GetDownloadURL(r.Context(), secret.Bucket, book.FilePath)

		userBook := database.UserBook{
			BookName:   book.Name,
			AuthorName: book.Name_2,
			ISBN:       book.Isbn,
		}

		stream := database.Stream{
			Book: userBook,
			Link: URL,
		}

		RespondWithJSON(w, http.StatusOK, stream)
	}
}
