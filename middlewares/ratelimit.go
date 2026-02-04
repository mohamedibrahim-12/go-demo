package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"go-demo/pkg/logger"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients     = make(map[string]*client)
	mu          sync.Mutex
	cleanupOnce sync.Once
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	c, ok := clients[ip]
	if !ok {
		// Default: 5 requests per second with burst of 10
		l := rate.NewLimiter(5, 10)
		clients[ip] = &client{limiter: l, lastSeen: time.Now()}
		return l
	}
	c.lastSeen = time.Now()
	return c.limiter
}

func startCleanup() {
	cleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				mu.Lock()
				for ip, c := range clients {
					if time.Since(c.lastSeen) > 5*time.Minute {
						delete(clients, ip)
					}
				}
				mu.Unlock()
			}
		}()
	})
}

// RateLimitMiddleware enforces a per-IP token-bucket limiter.
func RateLimitMiddleware(next http.Handler) http.Handler {
	startCleanup()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiter := getLimiter(ip)
		if !limiter.Allow() {
			logger.Log.Warn().Str("ip", ip).Msg("rate limit exceeded")
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
