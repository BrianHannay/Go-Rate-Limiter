/*
 * This package represents an example dependency included with your golang project.
 */
package main

// We need to output text and synchronize with parent processes
import (
	"time"

	"github.com/BrianHannay/Go-Rate-Limiter/ratelimit"
)

var New ratelimit.Constructor = func(attempts int, duration time.Duration) ratelimit.IRateLimit {
	return ratelimit.New(attempts, duration)
}
