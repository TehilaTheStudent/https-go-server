package main

import (
	"fmt"
	"net/http"
	"time"
)

// logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("--- Incoming Request ---")
		fmt.Printf("Time: %s\n", time.Now().Format(time.RFC3339))
		fmt.Printf("RemoteAddr: %s\n", r.RemoteAddr)
		fmt.Printf("Method: %s\n", r.Method)
		fmt.Printf("URL: %s\n", r.URL.String())
		fmt.Printf("Proto: %s\n", r.Proto)
		fmt.Println("Headers:")
		for k, v := range r.Header {
			fmt.Printf("  %s: %s\n", k, v)
		}
		if r.ContentLength > 0 {
			fmt.Printf("Content-Length: %d\n", r.ContentLength)
		}
		fmt.Println("-----------------------")
		next.ServeHTTP(w, r)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)

	loggedMux := loggingMiddleware(mux)

	fmt.Println("Server is running on https://0.0.0.0:8443")
	err := http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", loggedMux)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
