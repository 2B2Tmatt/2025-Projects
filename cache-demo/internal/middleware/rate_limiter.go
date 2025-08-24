package middleware

import (
	"log"
	"net/http"
	"time"
)

// use empty structs to send through channels to send signals without wasting too many resources
// important detail empty structs have size of 0 bytes, no overhead
type RateLimiter struct {
	tokens     chan struct{}
	refillTime time.Duration
}

func NewRateLimiter(rateLimit int, refillTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:     make(chan struct{}, rateLimit),
		refillTime: refillTime,
	}
	for range rateLimit {
		rl.tokens <- struct{}{}
	}
	go rl.Refill()
	return rl
}

func (rl *RateLimiter) Refill() {
	ticker := time.NewTicker(rl.refillTime)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}

	}
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func LimitRates(next http.Handler, rl *RateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.Allow() {
			log.Println("Request allowed")
			next.ServeHTTP(w, r)
		} else {
			log.Println("Rate limit exceeded")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	})
}
