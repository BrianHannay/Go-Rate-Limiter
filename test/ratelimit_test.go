// Packages prefixed with _test are run by "go test".
package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/BrianHannay/Go-Rate-Limiter/ratelimit"
)

// We'll need the testing package, for, you know... testing.

// Here's an example of testing a package using many subtests:
func TestDependency(t *testing.T) {

	// Map out all the tests we intend to run
	t.Run("RateLimit starts with capacity", testRateLimitDefaultsAllowed)
	t.Run("RateLimit limits rate", testLimiterLimits)
	t.Run("RateLimit runs synchronously", testSynchronousConsumption)

	t.Run("RateLimit race conditions", testLimitRace)
}

// Test that nil is an acceptable value to pass to dependency.Print
func testRateLimitDefaultsAllowed(t *testing.T) {
	ratelimiter := ratelimit.New(time.Second)
	if !ratelimiter.Allowed() {
		t.Errorf("RateLimit not allowed by default")
	}
	limit := ratelimiter.GetLimit()
	if limit != 1 {
		t.Errorf("RateLimit's limit does not match provided value, received: %d", ratelimiter.GetLimit())
	}
	attempts := ratelimiter.GetAttempts()
	if attempts != 0 {
		t.Errorf("RateLimit's attempts started at value other than 0, received: %d", ratelimiter.GetAttempts())
	}
	ratelimiter.Destroy()
}

// Test that Dependency.Print calls Done exactly once
func testLimiterLimits(t *testing.T) {
	ratelimiter := ratelimit.New(time.Minute)
	if !ratelimiter.ConsumeAsync() {
		t.Errorf("RateLimiter fails to consume an asynchronous attempt")
	}
	if ratelimiter.ConsumeAsync() {
		t.Errorf("RateLimiter fails to prevent requests")
	}
	ratelimiter.Destroy()
}

func testSynchronousConsumption(t *testing.T) {
	ratelimiter := ratelimit.New(time.Second)
	timeBeforeFirst := time.Now().UnixMicro()
	ratelimiter.Consume()
	timeBetweenAttempts := time.Now().UnixMicro()
	ratelimiter.Consume()
	timeAfter := time.Now().UnixMicro()
	ratelimiter.Destroy()

	if timeBetweenAttempts-timeBeforeFirst >= 1e6 {
		t.Errorf("First attempt took more than 1 second!")
	}
	if timeAfter-timeBeforeFirst < 1e6 {
		t.Errorf("Ratelimit took less than 1 second!")
	}
}

var consumes atomic.Int32
var group sync.WaitGroup

func countConsumes(ratelimiter ratelimit.IRateLimit) {
	// Loop for infinity until a request is allowed
	for {
		if ratelimiter.Allowed() {
			// Once allowed, all threads can call consumeAsync at once
			if ratelimiter.ConsumeAsync() {
				// And update the atomic counter if consume was successful
				consumes.Add(1)
			}
			group.Done()
			break
		} else if consumes.Load() > 0 {
			// If not allowed and something already incremented the atomic counter, that means this thread hit too late.
			// Bail out.
			group.Done()
			break
		}
	}
}

func testLimitRace(t *testing.T) {

	// 1 free consume per second
	ratelimiter := ratelimit.New(time.Second)

	// Clear the initial "free" consume
	ratelimiter.Consume()

	// Nothing has consumed yet
	consumes.Store(0)
	for range 10 {
		group.Add(1)
		// Try to consume all at once
		go countConsumes(ratelimiter)
	}
	// Wait for the threads to finish
	group.Wait()

	// Check how many succeeded in consuming
	if consumes.Load() != 1 {
		t.Errorf("Consumed an unexpected number of times: %d", consumes.Load())
	}
	ratelimiter.Destroy()
}
