package shared

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/opensaucerer/barf"
	domain "github.com/rafmme/anony-chat/pkg/domain/user"
	"golang.org/x/crypto/bcrypt"
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

func GenerateJWT(userData *domain.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userData.ID
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	env, err := LoadEnvVars()

	if err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	return token.SignedString([]byte(env.SecretKey))
}

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPasswordBytes), err
}

func PasswordMatches(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	return err == nil
}

func GetAuthToken(reqHeader http.Header, routeType string) string {
	var token string

	if routeType == "ws" {
		return strings.Replace(
			reqHeader.Get("Sec-Websocket-Protocol"),
			"Authorization, ", "", 1,
		)
	}

	token = strings.Replace(
		reqHeader.Get("Authorization"),
		"Bearer ", "", 1,
	)

	return token
}
