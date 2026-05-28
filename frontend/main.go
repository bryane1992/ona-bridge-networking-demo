package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	backendHost := os.Getenv("BACKEND_HOST")
	if backendHost == "" {
		backendHost = "backend"
	}
	backendURL := fmt.Sprintf("http://%s:8081/api/hello", backendHost)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get(backendURL)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, "failed to reach backend at %s: %v\n", backendURL, err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(w, "frontend got from backend: %s\n", string(body))
	})

	fmt.Printf("frontend listening on :8080, will call backend at %s\n", backendURL)
	http.ListenAndServe(":8080", nil)
}
