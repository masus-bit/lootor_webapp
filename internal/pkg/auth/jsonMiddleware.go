package auth

import "net/http"

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			w.Header().Set("X-Content-Type-Options", "nosniff")

			next.ServeHTTP(w, r)
		},
	)
}
