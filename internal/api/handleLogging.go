package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/knibirdgautam/library/internal/database"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("KEY")))

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   true, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	}
}

func HandleLogging(queries database.DBQueries) http.HandlerFunc {

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
		}

		hash, err := queries.GetPassHash(r.Context(), database.ToNullString(params.Email))
		if err != nil {
			RespondWithError(w, http.StatusBadGateway, "User Not Found")
			return
		}

		if database.PasswordVerify(hash, []byte(params.Password)) == false {
			return
		}

		user, err := queries.GetUser(r.Context(), database.ToNullString(params.Email))
		if err != nil {
			RespondWithError(w, http.StatusBadGateway, "Failed to get user")
			return
		}

		session, err := store.Get(r, "user-session")
		if err != nil {
			RespondWithError(w, http.StatusBadGateway, "Failed to get user")
			return
		}

		session.Values["Authenticated"] = true
		session.Values["User_id"] = user.ID

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
