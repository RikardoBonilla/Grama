package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type visitor struct {
	count     int
	windowEnd time.Time
}

// IPRateLimiter limits requests per client IP using an in-memory counter.
// It is safe for concurrent use. For multi-instance deployments, replace with
// a Redis-backed implementation.
type IPRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	window   time.Duration
}

func NewIPRateLimiter(rate int, window time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}
}

// Limit wraps next and returns 429 if the client IP exceeds the configured rate.
func (l *IPRateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)

		l.mu.Lock()
		now := time.Now()

		// Purge expired visitors on each request to avoid unbounded memory growth.
		for key, v := range l.visitors {
			if now.After(v.windowEnd) {
				delete(l.visitors, key)
			}
		}

		v, ok := l.visitors[ip]
		if !ok || now.After(v.windowEnd) {
			l.visitors[ip] = &visitor{count: 1, windowEnd: now.Add(l.window)}
			l.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		v.count++
		if v.count > l.rate {
			retryAfter := int(time.Until(v.windowEnd).Seconds()) + 1
			l.mu.Unlock()
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			http.Error(w, `{"error":"too many requests"}`, http.StatusTooManyRequests)
			return
		}
		l.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// extractIP returns the client IP, respecting X-Forwarded-For for reverse-proxied setups.
func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Use only the first (client) IP in the chain.
		if host, _, err := net.SplitHostPort(xff); err == nil {
			return host
		}
		return xff
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
