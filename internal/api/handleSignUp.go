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
			return
		}

		if params.FirstName == "" || params.LastName == "" || params.Email == "" || params.Password == "" {
			RespondWithError(w, http.StatusBadRequest, "First name, last name, email, and password are required")
			return
		}

		// Verify email uniqueness at application layer
		_, err = queries.GetUser(r.Context(), database.ToNullString(params.Email))
		if err == nil {
			RespondWithError(w, http.StatusConflict, "User with this email already exists")
			return
		}

		passHash, err := database.PasswordHash([]byte(params.Password))
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to process password")
			return
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
			RespondWithError(w, http.StatusBadRequest, "Failed to Create User: "+err.Error())
			return
		}

		err = queries.LinkHash(r.Context(), database.LinkHashParams{
			UserID:   user.ID,
			PassHash: passHash,
		})
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Unexpected Failure saving credentials")
			return
		}

		RespondWithJSON(w, http.StatusCreated, user)
	}
}
