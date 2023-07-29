package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/internal/interfaces/api/handlers"
	"github.com/rafmme/anony-chat/internal/interfaces/api/middleware"
	"github.com/rafmme/anony-chat/pkg/shared"
)

type Server struct {
	listenAddr string
}

func CreateServer() *Server {
	port := os.Getenv("PORT")

	if port == "" {
		log.Printf("PORT variable empty!")
	}

	return &Server{
		listenAddr: port,
	}
}

func (server *Server) Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Welcome to the one of a kind, Anonymous Messaging Platform.",
		})
	})

	http.Handle("/api/v1/signup", middleware.ValidateSignupData(
		handlers.SignupHandler,
	),
	)

	http.Handle("/api/v1/auth", middleware.ValidateAuthData(
		handlers.AuthHandler),
	)

	http.Handle("/ws/chat", middleware.Authenticate(
		handlers.HandleWebSocketConnection,
	),
	)

	defer shared.Database.Close()

	port := fmt.Sprintf(":%s", server.listenAddr)
	barf.Logger().Info(fmt.Sprintf("🆙 Server up on PORT %s", port))
	err := http.ListenAndServe(port, nil)

	if err != nil {
		barf.Logger().Error("Could'nt start the server. " + err.Error())
		os.Exit(1)
	}
}
