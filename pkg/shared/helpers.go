package shared

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/opensaucerer/barf"
	domain "github.com/rafmme/anony-chat/pkg/domain/user"
)

var (
	Database *gorm.DB
)

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

	Database = db

	// Migrate the schema
	db.AutoMigrate(&domain.User{})
}

func CreateUUID() string {
	return uuid.New().String()
}
