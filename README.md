# Golang Rate Limiter

A rate limiter allows restricting access to a system based on a frequency.

This package allows ratelimiting. Construct the limiter using `ratelimit.NewBursty(attempts int, duration time.Duration)`, then
consume the ratelimit using either `ratelimit.Consume()` for blocking consumption, or `ratelimit.ConsumeAsync() bool` for
non-blocking.

# Usage

A "bursty" rate limiter is applicable for most use-cases. This type of limiter allows a certain number of actions to be taken
before rate-limiting is applied. Construct a bursty rate limiter by with:
```go
var maximumRequests int = 5
var perInterval time.Duration = time.Minute
limiter := ratelimit.NewBursty(maximumRequests, perInterval)
```

You can also create a rate limiter that only allows a single action to be taken before another is allowed:
```go
var perInterval time.Duration = time.Second
limiter := ratelimit.New(perInterval)
```

After constructing the rate limiter, you can make an attempt against it as follows:
```go
if allowed := ratelimit.ConsumeAsync(); allowed {
    // your rate limited logic here
} else {
    // your error handler here
}
```

Alternatively, you can run code after waiting for a free attempt as follows:
```go
ratelimit.Consume()
// your ratelimited logic here
```

# Examples

## Refreshing the Limit
The limit is refreshed at a rate of `duration` / `attempts`. For example, if you've constructed a limiter and consumed all it's attempts:
```go
limiter := ratelimit.NewBursty(2, time.Minute)
limiter.Consume()
limiter.Consume()
```

Then the next attempt will be available in 30 seconds. After another 30 seconds, both attempts will be available again.

## Usage in a webserver
The main.go script is an example of how this could be used in an HTTP server. The source for the actual ratelimiter can
be found in the `ratelimiter/` directory.
