// main.go is the entry point for the web chat application.
package main

import (
	"log"
	"net/http"

	"web-chat-backend/handlers"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any source (for testing purposes; in production, it is better to restrict to a specific domain)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

	// Serve static files
	fs := http.FileServer(http.Dir("../chat-bot-frontend"))

	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.HandleFunc("/ws", handlers.HandleWebSocket)
	
	// Wrap the mux with CORS middleware
	http.Handle("/", enableCORS(mux))

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}

}
