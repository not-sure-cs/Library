package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
)

func HandleSignUp(queries database.DBQueries) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST requests allowed")
			return
		}

		type parameters struct {
			FirstName string `json:"firstname"`
			LastName  string `json:"lastname"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
			Password  string `json:"password"`
		}

		var params parameters

		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Couldn't Parse POST request")
		}

		user, err := queries.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FirstName: params.FirstName,
			LastName:  params.LastName,
			Email:     database.ToNullString(params.Email),
			PhNo:      database.ToNullString(params.Phone),
		})

		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to Create User")
		}

		err = queries.LinkHash(r.Context(), database.LinkHashParams{
			UserID:   user.ID,
			PassHash: database.PasswordHash([]byte(params.Password)),
		})

		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Unexpected Failure")
		}

		RespondWithJSON(w, http.StatusAccepted, user)
	}
}
