package api

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AuthedMiddleware(next http.Handler, store *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := store.Get(r, "user-session")
		if err != nil || session == nil || session.Values["Authenticated"] != true {
			RespondWithError(w, http.StatusUnauthorized, "Unauthorized: Valid session required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
