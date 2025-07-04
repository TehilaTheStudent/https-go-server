
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
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
	// Load server certificate and private key
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		fmt.Printf("Failed to load server cert/key: %v\n", err)
		return
	}

	// Load CA certificate to verify clients
	caCert, err := os.ReadFile("ca.pem")
	if err != nil {
		fmt.Printf("Failed to read CA cert: %v\n", err)
		return
	}
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		fmt.Println("Failed to append CA cert to pool")
		return
	}

	// Configure TLS with client certificate requirement
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	// Set up HTTP handler with middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	loggedMux := loggingMiddleware(mux)

	// Create and run HTTPS server
	server := &http.Server{
		Addr:      ":8443",
		Handler:   loggedMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("🔒 Server is running on https://0.0.0.0:8443 and requires client certificates")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
