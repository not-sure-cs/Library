package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

func HandleCreateBooks(queries *database.Queries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST requests allowed")
			return
		}

		type parameters struct {
			Title  string `json:"title"`
			Isbn   string `json:"isbn"`
			Author string `json:"author"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)

		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to decode POST request")
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
		})

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Could not Create Book")
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
