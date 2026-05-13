package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	count     int
	lastReset time.Time
}

type RateLimit struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
}

func NewRateLimit(limit int, window time.Duration) *RateLimit {
	return &RateLimit{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		rl.mu.Lock()

		c, exist := rl.clients[ip]

		if !exist {
			rl.clients[ip] = &client{
				count:     1,
				lastReset: time.Now(),
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(c.lastReset) > rl.window {
			c.count = 0
			c.lastReset = time.Now()
		}

		c.count++

		if c.count > rl.limit {
			rl.mu.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
