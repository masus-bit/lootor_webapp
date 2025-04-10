package auth

import (
	"context"
	"net/http"
	"strings"
)

func (s *JWTService) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		userLogin, err := s.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userLogin", userLogin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
