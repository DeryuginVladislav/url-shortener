package middleware

import (
	"context"
	"net/http"
	"strings"
	"url_shortener/configs"
	"url_shortener/pkg/jwt"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthed(w)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)

		if !isValid {
			writeUnauthed(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)

		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
