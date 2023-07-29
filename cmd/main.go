package main

import (
	_ "github.com/lib/pq"
	"github.com/rafmme/anony-chat/internal/interfaces/api"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func init() {
	shared.LoadEnvVars()
	shared.DbInitializers()
}

func main() {
	api.CreateServer().Start()
}
