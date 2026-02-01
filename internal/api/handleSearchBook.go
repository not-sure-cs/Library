package api

import (
	"lib/internal/database"
	"net/http"
	"strings"
)

func handleSearchBook(db *[]database.Book) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		author := query.Get("author")
		title := query.Get("title")

		results := *db

		if author == "" && title == "" {
			RespondWithError(w, http.StatusBadRequest, "All Fields Empty")
			return
		}

		if author != "" {
			var filtered []database.Book
			for _, b := range results {
				if strings.Contains(strings.ToLower(b.Author), strings.ToLower(author)) {
					filtered = append(filtered, b)
				}
			}
			results = filtered
		}

		if title != "" {
			var filtered []database.Book
			for _, b := range results {
				if strings.Contains(strings.ToLower(b.Title), strings.ToLower(title)) {
					filtered = append(filtered, b)
				}
			}
			results = filtered
		}

		RespondWithJSON(w, http.StatusOK, results)
	}
}
