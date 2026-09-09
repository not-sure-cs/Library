package api

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/knibirdgautam/library/internal/database"
)

func RequireRoles(store *sessions.CookieStore, RequiredRoles ...database.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "user-session")

			roleStr, ok := session.Values["user_role"].(string)
			if !ok {
				RespondWithError(w, http.StatusForbidden, "Forbidden: Missing role information")
				return
			}

			// 3. Verify if user role is in the allowed list
			for _, allowedRole := range RequiredRoles {
				if roleStr == string(allowedRole) {
					next.ServeHTTP(w, r)
					return
				}
			}
			RespondWithError(w, http.StatusForbidden, "Forbidden: Insufficient privileges")
		})
	}
}
