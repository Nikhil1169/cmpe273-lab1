package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Configure the client with a strict timeout
var client = &http.Client{Timeout: 1 * time.Second}

// GET /health
func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GET /call-echo?msg=...
func callEcho(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	msg := r.URL.Query().Get("msg")
	w.Header().Set("Content-Type", "application/json")

	// Call Service A
	url := fmt.Sprintf("http://127.0.0.1:8080/echo?msg=%s", msg)
	resp, err := client.Get(url)

	// FAILURE HANDLING: Service A is down or timed out
	if err != nil {
		log.Printf("service=B endpoint=/call-echo status=error error=%q latency_ms=%d", err.Error(), time.Since(start).Milliseconds())

		w.WriteHeader(http.StatusServiceUnavailable) // Return 503
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service_b": "ok",
			"service_a": "unavailable",
			"error":     err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// SUCCESS HANDLING
	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding response: %v", err)
	}

	log.Printf("service=B endpoint=/call-echo status=ok latency_ms=%d", time.Since(start).Milliseconds())
	_ = json.NewEncoder(w).Encode(map[string]any{
		"service_b": "ok",
		"service_a": data,
	})
}

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/call-echo", callEcho)

	log.Println("service=B listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
