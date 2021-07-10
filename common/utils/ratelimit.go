package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	lastResetTime time.Time
	rate          int64
	interval      time.Duration
	token         *int64
	mu            *sync.RWMutex
}

func NewRateLimiter(rate, intervalMillis int64) *RateLimiter {
	return &RateLimiter{
		lastResetTime: time.Now(),
		rate:          rate,
		interval:      time.Duration(intervalMillis) * time.Millisecond,
		token:         buildAddr(rate),
		mu:            new(sync.RWMutex),
	}
}

func (rl *RateLimiter) Allowable() bool {
	now := time.Now()
	if rl.isTimeOver(now) {
		rl.mu.Lock()
		rl.token = buildAddr(rl.rate)
		rl.lastResetTime = now
		rl.mu.Unlock()
	}

	rl.mu.RLock()
	if *rl.token <= 0 {
		return false
	}
	rl.decrement(rl.token)
	rl.mu.RUnlock()
	return true
}

func (rl *RateLimiter) isTimeOver(now time.Time) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.abs(now.UnixNano()-rl.lastResetTime.UnixNano()) > rl.interval.Nanoseconds()
}

func (rl *RateLimiter) decrement(ops *int64) {
	atomic.AddInt64(ops, -1)
}

func buildAddr(rate int64) *int64 {
	return &rate
}

func (rl *RateLimiter) abs(val int64) int64 {
	if val >= 0 {
		return val
	}
	return -val
}
