package api

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/knibirdgautam/library/internal/database"
)

func init() {
	gob.Register(uuid.UUID{})
	gob.Register(sql.NullString{})
}

func HandleLogging(queries database.DBQueries, store *sessions.CookieStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST requests allowed")
			return
		}

		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var params parameters

		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Couldn't Parse POST request")
			return
		}

		if params.Email == "" || params.Password == "" {
			RespondWithError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		hash, err := queries.GetPassHash(r.Context(), database.ToNullString(params.Email))
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		if !database.PasswordVerify(hash, []byte(params.Password)) {
			RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		user, err := queries.GetUser(r.Context(), database.ToNullString(params.Email))
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to get user")
			return
		}

		session, err := store.Get(r, "user-session")
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to get session")
			return
		}

		session.Values["Authenticated"] = true
		session.Values["User_id"] = user.ID
		session.Values["user_role"] = string(user.Role)
		session.Values["User_Role"] = string(user.Role)

		err = session.Save(r, w)
		if err != nil {
			log.Printf("Failed Session Creation :{%s}", err)
			RespondWithError(w, http.StatusInternalServerError, "Could not save session")
			return
		}

		log.Print("Successfully Created Session")
		RespondWithJSON(w, http.StatusOK, user)
	}

}
