/*
 * This package represents an example dependency included with your golang project.
 */
package ratelimit

// We need to output text and synchronize with parent processes
import (
	"sync"
	"time"
)

type rateLimit struct {
	mutex       sync.RWMutex
	channel     chan time.Time
	decayTicker *time.Ticker
	destroyed   chan bool
}

func New(duration time.Duration) IRateLimit {
	return NewBursty(1, duration)
}

func NewBursty(attempts int, duration time.Duration) IRateLimit {
	limiter := &rateLimit{
		channel:     make(chan time.Time, attempts),
		decayTicker: time.NewTicker(duration / time.Duration(attempts)),
		destroyed:   make(chan bool, 1),
	}
	go limiter.decayLoop()
	for range attempts {
		limiter.channel <- time.Now()
	}
	return limiter
}

func (limiter *rateLimit) Destroy() {
	limiter.decayTicker.Stop()
	limiter.destroyed <- true
	close(limiter.channel)
}

func (limiter *rateLimit) decayLoop() {
	for {
		select {
		case <-limiter.destroyed:
			close(limiter.destroyed)
			return
		case limiter.channel <- <-limiter.decayTicker.C:
		}
	}
}

func (limiter *rateLimit) GetLimit() int {
	return cap(limiter.channel)
}
func (limiter *rateLimit) GetAttempts() int {
	return limiter.GetLimit() - len(limiter.channel)
}

func (limiter *rateLimit) Allowed() bool {
	return len(limiter.channel) > 0
}

func (limiter *rateLimit) Consume() {
	limiter.mutex.Lock()
	<-limiter.channel
	limiter.mutex.Unlock()
}

func (limiter *rateLimit) ConsumeAsync() bool {
	if limiter.mutex.TryLock() {
		defer limiter.mutex.Unlock()
		if limiter.Allowed() {
			<-limiter.channel
			return true
		}
	}
	return false
}
