package middlewares

import "net/http"

// Hardcoded users
var allowedUsers = map[string]bool{
	"jennamandelricci": true,
	"asisatmuldoon":    true,
	"midrenelamy":      true,
	"adonisbrown":      true,
	"lindalaul":        true,
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Always allow login + logout routes
		if r.URL.Path == "/log-in" || r.URL.Path == "/sign-out" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("auth")
		if err != nil {
			http.Redirect(w, r, "/log-in", http.StatusSeeOther)
			return
		}

		// ✅ Direct lookup (no loop)
		if !allowedUsers[cookie.Value] {
			http.Redirect(w, r, "/log-in", http.StatusSeeOther)
			return
		}

		// ❗ Block "/" even if authenticated
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/talent-show", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}