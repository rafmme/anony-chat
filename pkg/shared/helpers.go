package shared

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/opensaucerer/barf"
	domain "github.com/rafmme/anony-chat/pkg/domain/user"
)

type Env struct {
	Port       string `barfenv:"key=PORT;required=true"`
	DbHost     string `barfenv:"key=DB_HOST;required=true"`
	DbPort     string `barfenv:"key=DB_PORT;required=true"`
	DbName     string `barfenv:"key=DB_NAME;required=true"`
	DbUser     string `barfenv:"key=DB_USER;required=true"`
	DbPassword string `barfenv:"key=DB_PASSWORD;required=true"`
}

func LoadEnvVars() (*Env, error) {

	env := new(Env)
	if err := barf.Env(env, ".env"); err != nil {
		return nil, err
	}

	return env, nil
}

func DbInitializers() {
	env, err := LoadEnvVars()

	if err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	var (
		dbHost     = env.DbHost
		dbPort     = env.DbPort
		dbName     = env.DbName
		dbUser     = env.DbUser
		dbPassword = env.DbPassword
	)

	databasePort, err := strconv.ParseInt(dbPort, 10, 32)
	if err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	dbURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		dbHost, databasePort, dbUser, dbName, dbPassword)

	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&domain.User{})
}
