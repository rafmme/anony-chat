package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/pkg/shared"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[string]*websocket.Conn)
var clientsLock sync.Mutex

func generateUniqueID() string {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	return fmt.Sprintf("client-%s", shared.CreateUUID()[0:8])
}

func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Generate a unique ID for the client
	clientID := generateUniqueID()

	// Store the WebSocket connection with the unique ID
	clientsLock.Lock()
	clients[clientID] = conn
	clientsLock.Unlock()

	log.Printf("Client connected with ID: %s", clientID)

	// Handle incoming messages from the client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			errMsg := fmt.Sprintf("Error reading message from client %s: %v", clientID, err)
			barf.Logger().Error(errMsg)
			break
		}

		log.Printf("Received message from client %s: %s", clientID, msg)

		// Example: Broadcast the message to all other clients
		clientsLock.Lock()
		for id, c := range clients {
			if id != clientID {
				c.WriteMessage(websocket.TextMessage, msg)
			}
		}
		clientsLock.Unlock()
	}

	// Remove the client from the list of connected clients when the connection is closed
	clientsLock.Lock()
	delete(clients, clientID)
	clientsLock.Unlock()

	log.Printf("Client %s disconnected", clientID)
}
