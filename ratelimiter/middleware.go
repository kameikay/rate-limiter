package ratelimiter

import (
	"context"
	"net"
	"net/http"
	"time"
)

type rateLimiterCheckFunction = func(ctx context.Context, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error)

func NewRateLimiter() func(next http.Handler) http.Handler {
	return NewRateLimiterConfig()
}

func NewRateLimiterConfig() func(next http.Handler) http.Handler {
	rateLimiterConfig := setConfig()
	return func(next http.Handler) http.Handler {
		return rateLimiter(rateLimiterConfig, next, checkRateLimit)
	}
}

func rateLimiter(config *RateLimiterConfig, next http.Handler, checkRateLimitFunc rateLimiterCheckFunction) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var blockTime *time.Time
		var err error

		token := r.Header.Get("API_KEY")
		if token != "" {
			tokenConfig := config.GetRateLimiterRateConfigForToken(token)
			blockTime, err = checkRateLimitFunc(r.Context(), token, config, tokenConfig)
		} else {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			blockTime, err = checkRateLimitFunc(r.Context(), host, config, config.IP)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if blockTime != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		}
	})
}
