package api

import (
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
			RespondWithError(w, http.StatusBadRequest, "Failed To Get Metadata")
			return
		}

		err = database.UnsaveFile(r.Context(), secret, store, metadata.FilePath)
		if err != nil {
			log.Printf("File deletion failed with:%v",err)
			RespondWithError(w, http.StatusBadRequest, "Failed To Execute Unsave")
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
