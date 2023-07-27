package main

import (
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rafmme/anony-chat/internal/interfaces/api"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func init() {
	shared.DbInitializers()
}

func main() {
	/* 	http.HandleFunc("/ws", handlers.HandleWebSocketConnection)
	   	log.Println("WebSocket server listening on :8080")
	   	log.Fatal(http.ListenAndServe(":8080", nil)) */
	api.StartServer()
}
