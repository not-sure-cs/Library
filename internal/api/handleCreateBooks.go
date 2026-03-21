package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/knibirdgautam/library/internal/database"
)

func HandleCreateBooks(q *database.Queries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST requests allowed")
			return
		}

		type parameters struct {
			Title  string `json:"title"`
			Isbn string `json:"isbn"`
			Author string `json:"author"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)

		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to decode POST request")
			return
		}

		aut, err := q.GetAuthor(r.Context(), params.Author)

		if err != nil {
			author, err := q.CreateAuthor(r.Context(), database.CreateAuthorParams{
				ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
				Name: paramsr,
			})

			if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Could not Create Author")
			
			}
		}


		book, err := q.CreateBook(r.Context(), database.CreateBookParams{
			ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
			CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			Name:     params.Title,
			Isbn: pgtype.Text{String: params.Isbn, Valid: true},

		})

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Could not Create Book")
			return 
		}

		linker, err:= q.LinkBookAuthor(r.Context(), database.LinkBookAuthorParams{
			BookID: book.ID ,
			AuthorID: author.ID,

		})

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Couldn't link Books and Authors")
		}

		RespondWithJSON(w, http.StatusOK, book)
	}
}
