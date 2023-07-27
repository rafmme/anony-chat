package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/opensaucerer/barf"
	infrastructure "github.com/rafmme/anony-chat/internal/infrastructure/persistence"
	"github.com/rafmme/anony-chat/internal/interfaces/api/middleware"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	existingUserData := infrastructure.UserRepo.FindByEmail(middleware.UserData.Email)

	if len(existingUserData.Email) > 0 {
		barf.Response(w).Status(http.StatusConflict).JSON(shared.ErrorResponse{
			StatusCode: 409,
			Errors: []map[string]string{
				{
					"email": fmt.Sprintf(
						"User with email address %s already exist on the app.", middleware.UserData.Email,
					),
				},
			},
			Message: "User already exist",
		})
		return
	}

	userData := infrastructure.UserRepo.Save(middleware.UserData)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userData.ID
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	env, err := shared.LoadEnvVars()

	if err != nil {
		barf.Logger().Error(err.Error())
		os.Exit(1)
	}

	tokenString, err := token.SignedString([]byte(env.SecretKey))
	if err != nil {
		barf.Logger().Error(err.Error())
		barf.Response(w).Status(http.StatusInternalServerError).JSON(shared.ErrorResponse{
			StatusCode: 500,
			Errors: []map[string]string{
				{
					"server": "Internal Server Error.",
				},
			},
			Message: "Internal Server Error.",
		})
		return
	}

	barf.Response(w).Status(http.StatusCreated).JSON(shared.Response{
		StatusCode: 201,
		Message:    "User signup was successful.",
		Data: map[string]map[string]string{
			"user": {
				"email":      userData.Email,
				"authToken":  tokenString,
				"_createdAt": userData.CreatedAt,
				"_updatedAt": userData.UpdatedAt,
			},
		},
	})
}
