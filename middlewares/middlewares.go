package middlewares

import "net/http"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow login page
		if r.URL.Path == "/log-in" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("auth")
		if err != nil || cookie.Value != "true" {
			http.Redirect(w, r, "/log-in", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}