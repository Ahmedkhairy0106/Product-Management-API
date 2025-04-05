package limiting

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiterStore struct {
	limiters map[string]*rate.Limiter
	mu sync.Mutex
	rate rate.Limit
	burst int
}

// function return ratelimit struct to create new rate limiter
func NewRateLimiterStore(rateLimit rate.Limit, burst int) *RateLimiterStore{
	// This function is used to create a new rate limiter store
	return &RateLimiterStore{
		limiters: make(map[string]*rate.Limiter),
		rate: rateLimit,
		burst: burst,
	}
}

// function to get rate limiter for a given ip address
func (s *RateLimiterStore) GetLimiter(ip string) *rate.Limiter{
	s.mu.Lock()
	defer s.mu.Unlock()
	// Check if the rate limiter for the IP address already exists
	limiter, exists := s.limiters[ip]

	// If it does not exist, create a new rate limiter and store it
	if !exists {
		limiter = rate.NewLimiter(s.rate, s.burst)
		s.limiters[ip] = limiter
	}
	// If it exists, return the existing rate limiter
	return limiter
}

// function to create rate limiter middleware
func RateLimiterMiddleware(store *RateLimiterStore) func(http.Handler) http.Handler{
	// This function is used to create a rate limiter middleware
	return func(next http.Handler) http.Handler {
		// This middleware function is used to limit the rate of requests
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the IP address from the request
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			// Get the rate limiter for the IP address
			limiter := store.GetLimiter(ip)

			// Check if the request is allowed
			if !limiter.Allow() {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			// Wait for the request to be allowed
			next.ServeHTTP(w, r)
		})
	}
}