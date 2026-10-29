package main

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type Limiter struct {
	mu     sync.Mutex
	clock  Clock
	rate   float64
	burst  int
	tokens float64
	last   time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	l := &Limiter{clock: clock, rate: ratePerSec, burst: burst}
	if burst > 0 && clock != nil {
		l.tokens = float64(burst)
		l.last = clock.Now()
	}
	return l
}
func (l *Limiter) AllowN(n int) bool {
	if l.burst <= 0 || l.clock == nil {
		return false
	}
	if n <= 0 || n > l.burst {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clock.Now()
	// A clock moved backwards must neither refill the bucket nor rewind last,
	// otherwise moving it forward again would refill the same interval twice.
	if elapsed := now.Sub(l.last); elapsed > 0 {
		l.last = now
		if l.rate > 0 {
			l.tokens += elapsed.Seconds() * l.rate
			if limit := float64(l.burst); l.tokens > limit {
				l.tokens = limit
			}
		}
	}
	// All or nothing: a short bucket keeps every token it has.
	if l.tokens < float64(n) {
		return false
	}
	l.tokens -= float64(n)
	return true
}
