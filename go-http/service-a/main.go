package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// GET /health
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /echo?msg=...
func echo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	msg := r.URL.Query().Get("msg")

	// Simulate processing time (optional, but good for testing timeouts later)
	// time.Sleep(50 * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"echo": msg})

	// Log the request details
	log.Printf("service=A endpoint=/echo status=ok latency_ms=%d", time.Since(start).Milliseconds())
}

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/echo", echo)

	log.Println("service=A listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
