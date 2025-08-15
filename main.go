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
	"time"

	"github.com/BrianHannay/golang-template-example/ratelimit"
)

var limiter *ratelimit.RateLimit

// the main function is the entrypoint to the compiled go program
func main() {

	requests := flag.Int("requests", 1, "Maximum number of requests per minute to handle")
	port := flag.Int("port", 8888, "Port number on which to listen")
	host := flag.String("host", "localhost", "Interface on which to listen")
	flag.Parse()

	addr := fmt.Sprintf(
		"%s:%d",
		*host,
		*port,
	)

	limiter = ratelimit.New(*requests, time.Minute)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	fmt.Println("Listening for any request")
	register(mux, "/sync", syncRequest)
	register(mux, "/async", asyncRequest)
	fmt.Println("Listening on", addr)
	server.ListenAndServe()
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
