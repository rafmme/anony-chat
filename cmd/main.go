package main

import (
	_ "github.com/dgrijalva/jwt-go"
	_ "github.com/google/uuid"
	_ "github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/rafmme/anony-chat/internal/interfaces/api"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func init() {
	shared.DbInitializers()
}

func main() {
	api.StartServer()
}
