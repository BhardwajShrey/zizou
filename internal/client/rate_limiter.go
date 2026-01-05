package client

import (
	"context"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	rate       int           // requests per period
	period     time.Duration // time period
	tokens     int           // available tokens
	lastRefill time.Time     // last refill time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, period time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		period:     period,
		tokens:     rate,
		lastRefill: time.Now(),
	}
}

// Wait waits until a token is available
func (rl *RateLimiter) Wait(ctx context.Context) error {
	for {
		rl.mu.Lock()

		// Refill tokens based on elapsed time
		rl.refill()

		// Check if we have tokens available
		if rl.tokens > 0 {
			rl.tokens--
			rl.mu.Unlock()
			return nil
		}

		// Calculate wait time
		waitTime := rl.period / time.Duration(rl.rate)
		rl.mu.Unlock()

		// Wait for next token
		select {
		case <-time.After(waitTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// refill refills tokens based on elapsed time
// Must be called with lock held
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	if elapsed >= rl.period {
		// Full refill
		rl.tokens = rl.rate
		rl.lastRefill = now
	} else {
		// Partial refill based on elapsed time
		tokensToAdd := int(float64(rl.rate) * (float64(elapsed) / float64(rl.period)))
		rl.tokens += tokensToAdd

		if rl.tokens > rl.rate {
			rl.tokens = rl.rate
		}

		if tokensToAdd > 0 {
			rl.lastRefill = now
		}
	}
}

// Available returns the number of available tokens
func (rl *RateLimiter) Available() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()
	return rl.tokens
}

// Reset resets the rate limiter
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.tokens = rl.rate
	rl.lastRefill = time.Now()
}
