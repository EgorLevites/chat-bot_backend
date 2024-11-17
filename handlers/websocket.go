package handlers

import (
	"log"
	"net/http"
	"web-chat-backend/models"
	"web-chat-backend/gemini"
	"github.com/gorilla/websocket"
)

var (
	
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// HandleWebSocket upgrades an HTTP connection to a WebSocket and manages communication with the client.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }
    defer ws.Close()

    log.Println("New client connected.")

    for {
        var msg models.Message
        err := ws.ReadJSON(&msg)
        if err != nil {
            log.Printf("JSON read error: %v", err)
            break
        }

        responseText, err := gemini.GenerateResponse(msg.Content)
        if err != nil {
            log.Println("Error generating response:", err)
            responseText = "Sorry, I couldn't generate a response."
        }

        err = ws.WriteJSON(models.Message{
            Username: "Gemini Bot",
            Content:  responseText,
        })
        if err != nil {
            log.Printf("Error writing to WebSocket: %v", err)
            break
        }
    }

    log.Println("Client disconnected.")
}


