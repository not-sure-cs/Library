package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
	"github.com/knibirdgautam/library/internal/storage"
)

func HandleCreateBooks(queries database.DBQueries, store storage.R2Store, secret storage.Secret) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST requests allowed")
			return
		}

		countChannel := make(chan int64, 1)
		errChannel := make(chan error, 1)

		go func() {
			num, err := queries.CountBook(r.Context())
			if err != nil {
				errChannel <- err
				return
			}
			countChannel <- num
		}()

		r.ParseMultipartForm(200 << 20)

		file, fileHandler, err := r.FormFile("uploadFile")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to Marshal File Response")
			return
		}

		var num int64
		select {
		case num = <-countChannel:

		case err = <-errChannel:
			RespondWithError(w, http.StatusBadRequest, "Couldn't Count Books")
			return
		}

		fileKey, err := database.SaveFile(num, r.Context(), secret, store, file, fileHandler)
		if err != nil {
			RespondWithError(w, 500, "Upload failed")
			return
		}

		type parameters struct {
			Title  string `json:"title"`
			Isbn   string `json:"isbn"`
			Author string `json:"author"`
		}
		jsonStr := r.FormValue("metadata")
		params := parameters{}

		err = json.NewDecoder(strings.NewReader(jsonStr)).Decode(&params)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to Decode JSON Body")
			return
		}

		var author database.Author
		author, err = queries.GetAuthor(r.Context(), params.Author)

		if err != nil {
			author, err = queries.CreateAuthor(r.Context(), database.CreateAuthorParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      params.Author,
			})

			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Could not Create Author")
				return
			}
		}

		book, err := queries.CreateBook(r.Context(), database.CreateBookParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      params.Title,
			Isbn:      database.ToNullString(params.Isbn),
			FilePath:  fileKey,
		})

		if err != nil {

			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("DBError: %s", err))
			return
		}

		linker, err := queries.LinkBookAuthor(r.Context(), database.LinkBookAuthorParams{
			BookID:   book.ID,
			AuthorID: author.ID,
		})

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Couldn't link Books and Authors")
			return
		}

		resp := database.Linked{
			Author: author,
			Book:   book,
			Link:   linker,
		}

		RespondWithJSON(w, http.StatusOK, resp)
	}
}
