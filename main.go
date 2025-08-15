/**
 * Template golang project
 *
 * This sample will show how WaitGroups can be used to ensure child threads have completed.
 * We will also see how channels can be used to communicate between threads.
 */
package main

// We are going to require some builtin features - namely, text formatting, and thread synchronization.
// Oh, let's also provide an example of using plugins!
import (
	"flag"
	"fmt"
	"net/http"
	"plugin"
	"time"

	"github.com/BrianHannay/Go-Rate-Limiter/ratelimit"
)

var limiter ratelimit.IRateLimit
var pluginLimiter ratelimit.IRateLimit

// the main function is the entrypoint to the compiled go program
func main() {

	requests := flag.Int("requests", 12, "Maximum number of requests per minute to handle")
	port := flag.Int("port", 8888, "Port number on which to listen")
	host := flag.String("host", "localhost", "Interface on which to listen")
	flag.Parse()

	addr := fmt.Sprintf(
		"%s:%d",
		*host,
		*port,
	)

	limiter = ratelimit.NewBursty(*requests, time.Minute)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	fmt.Println("Listening for any request")
	register(mux, "/sync", syncRequest)
	register(mux, "/async", asyncRequest)

	newRateLimiter := loadRatelimiterPlugin()
	if newRateLimiter != nil {
		pluginLimiter = newRateLimiter(*requests, time.Minute)
		register(mux, "/plugin/sync", pluginSyncRequest)
		register(mux, "/plugin/async", pluginAsyncRequest)
	}
	fmt.Println("Listening on", addr)
	server.ListenAndServe()
}

func loadRatelimiterPlugin() func(attempts int, duration time.Duration) ratelimit.IRateLimit {
	plugin, err := plugin.Open("./plugins/ratelimit.so")
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load plugin: %+v", err))
		return nil
	}

	symbol, err := plugin.Lookup("New")
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read RateLimit constructor: %+v", err))
		return nil
	}

	casted := symbol.(*ratelimit.BurstyConstructor)
	return *casted
}

func register(mux *http.ServeMux, route string, handleFunc func(http.ResponseWriter, *http.Request)) {
	fmt.Println("Listening for requests to", route)
	mux.HandleFunc(route, handleFunc)
}

func asyncRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to", r.RequestURI)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if limiter.ConsumeAsync() {
		time.Sleep(100 * time.Millisecond) // Some resource intensive operation
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
}

func syncRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to", r.RequestURI)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	limiter.Consume()
	time.Sleep(100 * time.Millisecond) // Some resource intensive operations
	w.Write([]byte("OK"))
}

func pluginAsyncRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to", r.RequestURI)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if pluginLimiter.ConsumeAsync() {
		time.Sleep(100 * time.Millisecond) // Some resource intensive operation
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
}

func pluginSyncRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to", r.RequestURI)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pluginLimiter.Consume()
	time.Sleep(100 * time.Millisecond) // Some resource intensive operations
	w.Write([]byte("OK"))
}
