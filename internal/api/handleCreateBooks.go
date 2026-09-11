package api

import (
	"encoding/json"
	"fmt"
	"log"
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

		err := r.ParseMultipartForm(200 << 20)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to parse multipart form")
			return
		}

		file, fileHandler, err := r.FormFile("uploadFile")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "uploadFile is required")
			return
		}
		defer file.Close()

		type parameters struct {
			Title        string `json:"title"`
			Isbn         string `json:"isbn"`
			Author       string `json:"author"`
			CategoryCode string `json:"category_code"`
			PubYear      int16  `json:"pub_year"`
		}
		jsonStr := r.FormValue("metadata")
		if jsonStr == "" {
			RespondWithError(w, http.StatusBadRequest, "Missing metadata form value")
			return
		}

		params := parameters{}
		err = json.NewDecoder(strings.NewReader(jsonStr)).Decode(&params)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to Decode JSON Body")
			return
		}

		if params.Title == "" || params.Author == "" {
			RespondWithError(w, http.StatusBadRequest, "Title and Author are required")
			return
		}

		fileKey, err := database.SaveFile(0, r.Context(), secret, store, file, fileHandler)
		if err != nil {
			log.Printf("R2 Upload Error: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Upload failed")
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
				_ = database.UnsaveFile(r.Context(), secret, store, fileKey)
				RespondWithError(w, http.StatusInternalServerError, "Could not Create Author")
				return
			}
		}

		book, err := queries.CreateBook(r.Context(), database.CreateBookParams{
			ID:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Name:         params.Title,
			Isbn:         database.ToNullString(params.Isbn),
			FilePath:     fileKey,
			CategoryCode: params.CategoryCode,
			PubYear:      params.PubYear,
		})

		if err != nil {
			_ = database.UnsaveFile(r.Context(), secret, store, fileKey)
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("DBError: %s", err))
			return
		}

		linker, err := queries.LinkBookAuthor(r.Context(), database.LinkBookAuthorParams{
			BookID:   book.ID,
			AuthorID: author.ID,
		})

		if err != nil {
			_ = database.UnsaveFile(r.Context(), secret, store, fileKey)
			RespondWithError(w, http.StatusInternalServerError, "Couldn't link Books and Authors")
			return
		}

		resp := database.Linked{
			Author: author,
			Book:   book,
			Link:   linker,
		}

		RespondWithJSON(w, http.StatusCreated, resp)
	}
}
