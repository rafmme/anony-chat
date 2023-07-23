package api

import (
	"net/http"
	"os"

	"github.com/opensaucerer/barf"
)

func StartServer() {
	type Env struct {
		Port string `barfenv:"key=PORT;required=true"`
	}

	env := new(Env)
	if err := barf.Env(env, ".env"); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
	allow := true

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

	if err := barf.Beck(); err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}
