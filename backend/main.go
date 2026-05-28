package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"message": "hello from backend"}`)
	})
	fmt.Println("backend listening on :8081")
	http.ListenAndServe(":8081", nil)
}
