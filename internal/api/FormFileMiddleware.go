package api

import "net/http"

func FormFileMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "multipart/form")

		next.ServeHTTP(w, r)
	})
}
