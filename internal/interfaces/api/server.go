package api

import (
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/internal/interfaces/api/handlers"
	"github.com/rafmme/anony-chat/internal/interfaces/api/middleware"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func StartServer() {
	allow := true
	env, err := shared.LoadEnvVars()

	if err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	if err := barf.Stark(barf.Augment{
		Port:     env.Port,
		Logging:  &allow,
		Recovery: &allow,
	}); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	barf.Get("/", func(w http.ResponseWriter, r *http.Request) {
		barf.Response(w).Status(http.StatusOK).JSON(barf.Res{
			Status:  true,
			Data:    nil,
			Message: "Welcome to the one of a kind, Anonymous Messaging Platform.",
		})
	})

	apiRouter := barf.RetroFrame("/api").RetroFrame("/v1")
	apiRouter.Post("/signup", middleware.ValidateSignupData(
		handlers.SignupHandler,
	),
	)

	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
