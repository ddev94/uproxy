package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// Define the target URL to proxy to
	targetURL := flag.String("target-url", "", "Target URL to proxy to (e.g., https://google.com)")
	flag.Parse()
	fmt.Println(*targetURL)

	// Create a handler for proxying requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the target URL
		target, err := url.Parse(*targetURL)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			log.Printf("Error parsing target URL: %v", err)
			return
		}

		// Create a new request to the target
		proxyUrl := target.String() + r.URL.Path
		if r.URL.RawQuery != "" {
			proxyUrl += "?" + r.URL.RawQuery
		}
		proxyReq, err := http.NewRequest(r.Method, proxyUrl, r.Body)
		if err != nil {
			http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
			log.Printf("Error creating proxy request: %v", err)
			return
		}

		// Copy headers from original request
		for header, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(header, value)
			}
		}

		// Create HTTP client and send the request
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, "Error forwarding request", http.StatusBadGateway)
			log.Printf("Error forwarding request: %v", err)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for header, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(header, value)
			}
		}

		// Set status code
		w.WriteHeader(resp.StatusCode)

		// Copy response body
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("Error copying response body: %v", err)
		}
	})

	// Start the server
	port := 8080
	log.Printf("Starting proxy server on :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
