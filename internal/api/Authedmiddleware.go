package api

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func AuthedMiddleware(next http.Handler, store *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "user-session")
		if session.Values["Authenticated"] == true {
			next.ServeHTTP(w, r)
		}
	})
}
