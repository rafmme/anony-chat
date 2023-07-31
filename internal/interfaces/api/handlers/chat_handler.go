package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/pkg/shared"
)

const serverId string = "ac"

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
	return fmt.Sprintf("user-%s", userId[0:8])
}

func sendClientCount() {
	clientsList := map[string]string{
		serverId: serverId,
	}

	for clientId := range clients {
		clientsList[clientId] = clientId
	}

	sendServerMessage(&shared.Message{
		MsgType:     "count",
		ClientCount: len(clientsList),
		ClientsList: clientsList,
	})
}

func sendServerMessage(message *shared.Message) {
	for id, client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Println("Error sending client count:", err)
			client.Close()
			delete(clients, id)
		}
	}
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

		sendClientCount()
		sendServerMessage(&shared.Message{
			MsgType:  "msg",
			Action:   "joined",
			ClientID: clientID,
			Message:  fmt.Sprintf("%s has joined the chat.", clientID),
			Sender:   serverId,
			Date:     time.Now(),
		})
		log.Printf("Client connected with ID: %s", clientID)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			errMsg := fmt.Sprintf("Error reading message from client %s: %v", clientID, err)
			barf.Logger().Error(errMsg)
			break
		}

		text := string(msg)
		log.Printf("Received message from %s: %s", clientID, msg)

		mentionList := shared.GetUserIDsInChatMessage(text, "@")
		privateMessageList := shared.GetUserIDsInChatMessage(text, ".")

		// Example: Broadcast the message to all other clients
		clientsLock.Lock()

		if len(privateMessageList) > 0 {
			for id, c := range clients {
				if shared.CheckIfStringInSlice(privateMessageList, id) || id == clientID {
					err := c.WriteJSON(map[string]interface{}{
						"msgType": "msg",
						"message": shared.RemoveUsersIDFromMessage(text),
						"sender":  clientID,
						"private": true,
						"date":    time.Now(),
					})

					if err != nil {
						log.Println("Error sending client message:", err)
						c.Close()
						delete(clients, id)
					}
				}
			}
		} else {
			for id, c := range clients {
				var err error
				if shared.CheckIfStringInSlice(mentionList, id) {
					err = c.WriteJSON(map[string]interface{}{
						"msgType":   "msg",
						"message":   text,
						"sender":    clientID,
						"mentioned": true,
						"date":      time.Now(),
					})
				} else {
					err = c.WriteJSON(map[string]interface{}{
						"msgType": "msg",
						"message": text,
						"sender":  clientID,
						"date":    time.Now(),
					})
				}

				if err != nil {
					log.Println("Error sending client message:", err)
					c.Close()
					delete(clients, id)
				}
			}
		}

		clientsLock.Unlock()
	}

	// Remove the client from the list of connected clients when the connection is closed
	clientsLock.Lock()
	delete(clients, clientID)
	clientsLock.Unlock()

	sendClientCount()
	sendServerMessage(&shared.Message{
		MsgType: "msg",
		Message: fmt.Sprintf("%s has left the chat.", clientID),
		Sender:  serverId,
		Date:    time.Now(),
	})
	log.Printf("Client %s disconnected", clientID)
}
