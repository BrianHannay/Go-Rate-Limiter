# Golang Rate Limiter

This package allows ratelimiting. Construct the limiter using `ratelimit.New(attempts int, duration time.Duration)`, then
consume the ratelimit using either `ratelimit.Consume()` for blocking consumption, or `ratelimit.ConsumeAsync() bool` for
non-blocking.

The main.go script is just an example of how this could be used in an HTTP server. The actual ratelimiter is in `ratelimiter/`.