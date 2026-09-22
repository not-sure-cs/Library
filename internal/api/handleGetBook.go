package api

import (
	"database/sql"
	"errors"
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
			if errors.Is(err, sql.ErrNoRows) {
				RespondWithError(w, http.StatusNotFound, "Book not found")
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}

		URL, err := store.GetDownloadURL(r.Context(), secret.Bucket, book.FilePath)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to generate download URL")
			return
		}

		userBook := database.UserBook{
			ID:         book.ID,
			BookName:   book.BookName,
			AuthorName: book.AuthorName,
			ISBN:       book.Isbn,
			MimeType:   book.MimeType,
			PageCount:  book.PageCount,
			Producer:   book.Producer,
			Subject:    book.Subject,
			PdfVersion: book.PdfVersion,
			CoverPath:  book.CoverPath,
			CreatedAt:  book.CreatedAt,
			UpdatedAt:  book.UpdatedAt,
		}

		stream := database.Stream{
			Book: userBook,
			Link: URL,
		}

		RespondWithJSON(w, http.StatusOK, stream)
	}
}
