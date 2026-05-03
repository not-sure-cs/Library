package api

import (
	"net/http"
)

func AuthedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := store.Get(r, "user-session")
		if session.Values["Authenticated"] == true {
			next.ServeHTTP(w, r)
		}
	})
}
