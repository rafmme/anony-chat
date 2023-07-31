package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/internal/interfaces/api/handlers"
	"github.com/rafmme/anony-chat/internal/interfaces/api/middleware"
	"github.com/rafmme/anony-chat/pkg/shared"
)

type Server struct {
	listenAddr string
}

func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)
		log.Printf("%s %s - %v", r.Method, r.URL.Path, duration)
	})
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
	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()

	mux.Handle("/", fs)

	mux.Handle("/api/v1/auth/signup", middleware.ValidateSignupData(
		handlers.SignupHandler,
	),
	)

	mux.Handle("/api/v1/auth/login", middleware.ValidateAuthData(
		handlers.AuthHandler,
	),
	)

	mux.Handle("/ws/chat", middleware.Authenticate(
		handlers.HandleWebSocketConnection,
	),
	)

	defer shared.Database.Close()

	loggedMux := requestLoggerMiddleware(mux)
	port := fmt.Sprintf(":%s", server.listenAddr)
	barf.Logger().Info(fmt.Sprintf("🆙 Server up on PORT %s", port))
	err := http.ListenAndServe(port, loggedMux)

	if err != nil {
		barf.Logger().Error("Could'nt start the server. " + err.Error())
		os.Exit(1)
	}
}
