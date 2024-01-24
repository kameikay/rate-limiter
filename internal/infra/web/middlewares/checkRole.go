package middlewares

import (
	"net/http"
)

func RateLimiter(roles []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			return
			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}
