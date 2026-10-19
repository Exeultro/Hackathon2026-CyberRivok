package httpserver

import (
	"net/http"
	"sync"
	"time"

	"typing/internal/httpx"
)

type rlBucket struct {
	count int
	reset time.Time
}

func rateLimitAuth(limit int, window time.Duration) func(http.Handler) http.Handler {
	var (
		mu sync.Mutex
		m  = map[string]*rlBucket{}
	)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			now := time.Now()

			mu.Lock()
			b := m[ip]
			if b == nil || now.After(b.reset) {
				b = &rlBucket{count: 0, reset: now.Add(window)}
				m[ip] = b
			}
			b.count++
			exceeded := b.count > limit
			mu.Unlock()

			if exceeded {
				httpx.WriteError(w, r, httpx.ErrRateLimited)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
