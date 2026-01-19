package middleware

import (
	"io"
	"net/http"
	"strings"
)

func RequireMethod(method string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireContentType(rct string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ct := r.Header.Get("Content-Type")
			if !strings.HasPrefix(ct, rct) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireNonEmptyBody() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n, err := io.ReadAll(r.Body)
			if err != nil && err.Error() != "EOF" || len(n) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(strings.NewReader(string(n)))
			next.ServeHTTP(w, r)
		})
	}
}
