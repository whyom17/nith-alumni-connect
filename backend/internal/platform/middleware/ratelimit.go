package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nith-alumni-connect/backend/internal/platform/httputil"
)

type clientRecord struct {
	tokens     int
	lastRefill time.Time
}

// RateLimiter implements an in-memory token bucket rate limiter per client IP.
type RateLimiter struct {
	mu         sync.Mutex
	clients    map[string]*clientRecord
	maxTokens  int
	refillRate time.Duration // e.g. 1 token per duration
}

// NewRateLimiter creates a RateLimiter allowing maxTokens burst, refilling 1 token every refillRate.
func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:    make(map[string]*clientRecord),
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}

	// Periodic cleanup of stale client records every 10 minutes
	go rl.cleanupRoutine()
	return rl
}

func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, record := range rl.clients {
			if now.Sub(record.lastRefill) > 30*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Limit returns a middleware that rate limits requests based on client IP.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)

		rl.mu.Lock()
		record, exists := rl.clients[ip]
		now := time.Now()

		if !exists {
			record = &clientRecord{
				tokens:     rl.maxTokens - 1,
				lastRefill: now,
			}
			rl.clients[ip] = record
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Refill tokens based on elapsed time
		elapsed := now.Sub(record.lastRefill)
		tokensToAdd := int(elapsed / rl.refillRate)
		if tokensToAdd > 0 {
			record.tokens += tokensToAdd
			if record.tokens > rl.maxTokens {
				record.tokens = rl.maxTokens
			}
			record.lastRefill = now
		}

		if record.tokens <= 0 {
			rl.mu.Unlock()
			w.Header().Set("Retry-After", "60")
			httputil.WriteError(w, http.StatusTooManyRequests, "rate_limit_exceeded", "too many requests, please try again later", nil)
			return
		}

		record.tokens--
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func extractIP(r *http.Request) string {
	// If X-Forwarded-For is present, take the first IP
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
