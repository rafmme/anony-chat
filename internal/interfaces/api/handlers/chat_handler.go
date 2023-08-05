package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
	return fmt.Sprintf("user-%s", userId[:8])
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

	action := "joined"

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

	clientIP, err := shared.GetClientIP(r)
	if err != nil {
		log.Println("No Client IP: ", err.Error())
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err, userID)
		return
	}
	defer conn.Close()

	clientID := generateUserID(userID)

	clientsLock.Lock()

	if clients[clientID] != nil {
		if err := conn.Close(); err != nil {
			barf.Logger().Error(err.Error())
		}

		action = "rejoined"
		conn = clients[clientID]
	} else {
		clients[clientID] = conn
	}

	clientsLock.Unlock()

	sendClientCount()
	sendServerMessage(&shared.Message{
		MsgType:  "msg",
		Action:   action,
		ClientID: clientID,
		Message:  fmt.Sprintf("@%s has %s the chat.", clientID, action),
		Sender:   serverId,
		Date:     time.Now(),
	})
	log.Printf("Client connected with ID: %s and IP Address: %s", clientID, clientIP)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			errMsg := fmt.Sprintf("Error reading message from client %s: %v", clientID, err)
			barf.Logger().Error(errMsg)
			break
		}

		text := string(msg)
		log.Printf("Received message from %s: %s", clientID, msg)

		clientsLock.Lock()
		processClientMessage(clientID, text)
		clientsLock.Unlock()
	}

	// Remove the client from the list of connected clients when the connection is closed
	clientsLock.Lock()
	delete(clients, clientID)
	clientsLock.Unlock()

	sendClientCount()
	sendServerMessage(&shared.Message{
		MsgType: "msg",
		Message: fmt.Sprintf("@%s has left the chat.", clientID),
		Sender:  serverId,
		Date:    time.Now(),
	})
	log.Printf("Client %s disconnected", clientID)
}

func processClientMessage(clientID, text string) {
	if strings.ToLower(text) == ".q" || strings.ToLower(text) == ".q!" {
		err := clients[clientID].Close()
		if err != nil {
			barf.Logger().Error(err.Error())
			return
		}

		delete(clients, clientID)
		return
	}

	mentionList := shared.GetUserIDsInChatMessage(text, "@")
	privateMessageList := shared.GetUserIDsInChatMessage(text, ".")

	if (len(mentionList) == 1 &&
		strings.ToLower(mentionList[0]) == serverId) ||
		(len(privateMessageList) == 1 &&
			strings.ToLower(privateMessageList[0]) == serverId) {

		if strings.HasPrefix(text, "@") {
			sendServerMessage(&shared.Message{
				MsgType: "msg",
				Message: text,
				Sender:  clientID,
				Date:    time.Now(),
			})
		} else {
			err := clients[clientID].WriteJSON(&shared.Message{
				MsgType: "msg",
				Message: text,
				Sender:  clientID,
				Private: true,
				To:      []string{serverId},
				Date:    time.Now(),
			})

			if err != nil {
				log.Println("Error sending client message:", err)
				clients[clientID].Close()
				delete(clients, clientID)
			}
		}

		botResponse := shared.ServerChat()
		if len(botResponse) > 0 {
			if strings.HasPrefix(text, "@") {
				sendServerMessage(&shared.Message{
					MsgType:   "msg",
					Message:   fmt.Sprintf("@%s: %s", clientID, botResponse),
					Sender:    serverId,
					Mentioned: true,
					Date:      time.Now(),
				})
			} else {
				err := clients[clientID].WriteJSON(&shared.Message{
					MsgType: "msg",
					Message: botResponse,
					Sender:  serverId,
					Private: true,
					To:      []string{clientID},
					Date:    time.Now(),
				})

				if err != nil {
					log.Println("Error sending client message:", err)
					clients[clientID].Close()
					delete(clients, clientID)
				}
			}
		}
	} else {
		if len(privateMessageList) > 0 {
			for id, c := range clients {
				if shared.CheckIfStringInSlice(privateMessageList, id) || id == clientID {
					err := c.WriteJSON(&shared.Message{
						MsgType: "msg",
						Message: shared.RemoveUsersIDFromMessage(text),
						Sender:  clientID,
						To:      privateMessageList,
						Private: true,
						Date:    time.Now(),
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
					err = c.WriteJSON(&shared.Message{
						MsgType:   "msg",
						Message:   text,
						Sender:    clientID,
						Mentioned: true,
						Date:      time.Now(),
					})
				} else {
					err = c.WriteJSON(&shared.Message{
						MsgType: "msg",
						Message: text,
						Sender:  clientID,
						Date:    time.Now(),
					})
				}

				if err != nil {
					log.Println("Error sending client message:", err)
					c.Close()
					delete(clients, id)
				}
			}
		}
	}
}
