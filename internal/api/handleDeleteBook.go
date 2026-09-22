package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
	"github.com/knibirdgautam/library/internal/storage"
)

func HandleDeleteBook(queries database.DBQueries, store storage.R2Store, secret storage.Secret) http.HandlerFunc {

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

		metadata, err := queries.GetMetaData(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondWithError(w, http.StatusNotFound, "Book not found")
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}

		err = queries.UnlinkBook(r.Context(), id)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Couldn't Unlink book")
			return
		}

		err = queries.DeleteBook(r.Context(), id)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Couldn't Delete book")
			return
		}

		if metadata.FilePath != "" {
			err1 := database.UnsaveFile(r.Context(), secret, store, metadata.FilePath)
			err2 := database.UnsaveFile(r.Context(), secret, store, metadata.CoverPath)
			if err1 != nil || err2 != nil {
				log.Printf("File deletion from storage failed with: %v & %v", err1, err2)
			}
		}

		RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Book deleted"})
	}
}
