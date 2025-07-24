package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	targetURL := flag.String("target-url", "", "Target URL to proxy to (e.g., https://google.com)")
	flag.Parse()
	if *targetURL == "" {
		log.Fatal("Target URL is required")
	}
	fmt.Println("Proxying to:", *targetURL)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		target, err := url.Parse(*targetURL)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusInternalServerError)
			log.Printf("Error parsing target URL: %v", err)
			return
		}

		// Create a new request to the target
		proxyURL := target.String() + r.URL.Path
		if r.URL.RawQuery != "" {
			proxyURL += "?" + r.URL.RawQuery
		}
		req, err := http.NewRequest(r.Method, proxyURL, r.Body)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("client: got response!\n")
		fmt.Printf("client: status code: %d\n", res.StatusCode)

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours
		w.WriteHeader(res.StatusCode)
		w.Write(resBody)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
