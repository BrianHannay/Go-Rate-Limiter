package ratelimit

import "time"

type IRateLimit interface {
	Destroy()
	decayLoop()
	GetLimit() int
	GetAttempts() int
	Allowed() bool
	Consume()
	ConsumeAsync() bool
}

type Constructor func(attempts int, duration time.Duration) IRateLimit
