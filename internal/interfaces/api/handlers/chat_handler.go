package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
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

func generateUserID(userId string) string {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	return fmt.Sprintf("client-%s", userId[0:8])
}

func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(shared.AuthData{}).(jwt.MapClaims)
	userID, ok := claims["sub"].(string)

	if !ok {
		barf.Response(w).Status(http.StatusUnauthorized).JSON(shared.ErrorResponse{
			StatusCode: 401,
			Errors: []map[string]string{
				{
					"auth": "Unauthorized.",
				},
			},
			Message: "Unauthorized.",
		})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err, userID)
		return
	}
	defer conn.Close()

	clientID := generateUserID(userID)

	if clients[clientID] == nil {
		clientsLock.Lock()
		clients[clientID] = conn
		clientsLock.Unlock()

		log.Printf("Client connected with ID: %s", clientID)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			errMsg := fmt.Sprintf("Error reading message from client %s: %v", clientID, err)
			barf.Logger().Error(errMsg)
			break
		}

		log.Printf("Received message from %s: %s", clientID, msg)

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
