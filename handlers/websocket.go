package handlers

import (
	"log"
	"net/http"
	"web-chat-backend/models" // Importing the Message model to structure messages
	"web-chat-backend/gemini" // Importing Gemini package for AI response generation
	"github.com/gorilla/websocket" // WebSocket package for handling WebSocket connections
)

var (
	// Configuring WebSocket upgrader to upgrade HTTP connections to WebSocket
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024, // Setting buffer size for reading WebSocket messages
		WriteBufferSize: 1024, // Setting buffer size for writing WebSocket messages
		// Allowing all origins for simplicity; adjust for stricter security in production
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// HandleWebSocket upgrades an HTTP connection to a WebSocket and manages communication with the client.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close() // Ensure the WebSocket connection is closed when the function ends

	log.Println("New client connected.") // Log when a new client connects

	for {
		var msg models.Message // Create a variable to hold the incoming message

		// Read and parse JSON message from the WebSocket
		err := ws.ReadJSON(&msg)
		if err != nil {
			// Log the error if message reading fails and break the loop to disconnect the client
			log.Printf("JSON read error: %v", err)
			break
		}

		// Generate a response for the incoming message using the Gemini API
		responseText, err := gemini.GenerateResponse(msg.Content)
		if err != nil {
			// Log the error and set a default response if the AI fails to generate a response
			log.Println("Error generating response:", err)
			responseText = "Sorry, I couldn't generate a response."
		}

		// Send the generated response back to the client as a JSON message
		err = ws.WriteJSON(models.Message{
			Username: "Gemini Bot", // Set the username as "Gemini Bot" for the response
			Content:  responseText, // Set the response content
		})
		if err != nil {
			// Log the error if sending the message fails and break the loop to disconnect the client
			log.Printf("Error writing to WebSocket: %v", err)
			break
		}
	}

	log.Println("Client disconnected.") // Log when the client disconnects
}
