package ratelimit

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	limits   map[uint64]*TokenBucket
	capacity int
	rate     time.Duration
}

type TokenBucket struct {
	tokens     int
	lastRefill time.Time
	capacity   int
	rate       time.Duration
}

func NewRateLimiter(capacity int, rate time.Duration) *RateLimiter {
	return &RateLimiter{
		limits:   make(map[uint64]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}
}

func (r1 *RateLimiter) Allow(userID *uint64) bool {
	r1.mu.Lock()
	defer r1.mu.Unlock()

	id := uint64(0)
	if userID != nil {
		id = *userID
	}

	bucket, exists := r1.limits[id]
	if !exists {
		bucket = &TokenBucket{
			tokens:     r1.capacity,
			lastRefill: time.Now(),
			capacity:   r1.capacity,
			rate:       r1.rate,
		}
		r1.limits[id] = bucket
	}

	return bucket.consume()
}

func (tb *TokenBucket) consume() bool {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed / tb.rate)

	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
