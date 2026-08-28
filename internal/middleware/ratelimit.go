package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimit allows a fixed number of requests per client IP within a sliding window.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 60
	}
	if window <= 0 {
		window = time.Minute
	}

	store := struct {
		sync.Mutex
		clients map[string][]time.Time
	}{
		clients: make(map[string][]time.Time),
	}

	return func(c *gin.Context) {
		key := c.ClientIP()
		now := time.Now()
		cutoff := now.Add(-window)

		store.Lock()
		entries := store.clients[key]
		filtered := entries[:0]
		for _, ts := range entries {
			if ts.After(cutoff) {
				filtered = append(filtered, ts)
			}
		}

		if len(filtered) >= limit {
			store.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		filtered = append(filtered, now)
		store.clients[key] = filtered
		store.Unlock()

		c.Next()
	}
}
