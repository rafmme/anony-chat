package shared

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/opensaucerer/barf"
	domain "github.com/rafmme/anony-chat/pkg/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	Database *gorm.DB
)

func LoadEnvVars() {
	err := godotenv.Load()

	if err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}
}

func DbInitializers() {
	var (
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
		dbName     = os.Getenv("DB_NAME")
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
	)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		barf.Logger().Error("Failed to connect to database: " + err.Error())
		os.Exit(1)
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

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
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

func GetAuthToken(r *http.Request, routeType string) string {
	var token string

	if routeType == "ws" {
		s, err := url.QueryUnescape(r.URL.Query().Get("auth_token"))

		if err != nil {
			return ""
		}

		return s
	}

	token = strings.Replace(
		r.Header.Get("Authorization"),
		"Bearer ", "", 1,
	)

	return token
}

func GetUserIDsInChatMessage(chat, prefix string) []string {
	userIDs := []string{}

	if strings.HasPrefix(chat, prefix) {
		if strings.HasPrefix(chat, prefix+"[") {
			formattedChat := strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						strings.Split(chat, " ")[0], prefix, ""),
					"[", ""), "]", "",
			)

			userIDs = strings.Split(formattedChat, ",")
			return userIDs
		}

		userIDs = append(userIDs,
			strings.ReplaceAll(
				strings.Split(chat, " ")[0],
				prefix, ""),
		)

		return userIDs
	}

	return userIDs
}

func CheckIfStringInSlice(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func RemoveUsersIDFromMessage(str string) string {
	chatArray := strings.Split(str, " ")
	chatArray[0] = ""
	return strings.Trim(strings.Join(chatArray, " "), " ")
}

func GenerateRandomNumber(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
