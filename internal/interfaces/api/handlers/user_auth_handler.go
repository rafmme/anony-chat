package handlers

import (
	"fmt"
	"net/http"

	"github.com/opensaucerer/barf"
	application "github.com/rafmme/anony-chat/internal/app/user"
	infrastructure "github.com/rafmme/anony-chat/internal/infrastructure/persistence"
	"github.com/rafmme/anony-chat/pkg/shared"
)

var (
	userService *application.UserService
)

func init() {
	userService = &application.UserService{
		UserRepository: &infrastructure.UserRepository{},
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(shared.UserSignupData{}).(*shared.UserSignupData)

	if userData == nil {
		barf.Logger().Error("no user data")
		barf.Response(w).Status(http.StatusBadRequest).JSON(shared.ErrorResponse{
			StatusCode: 400,
			Errors: []map[string]string{
				{
					"email": "Requires `email` field.",
				},
				{
					"password": "Requires `password` field.",
				},
				{
					"confirmPassword": "Requires `confirmPassword` field.",
				},
			},
			Message: "Invalid request body.",
		})
		return
	}

	existingUserData, err := userService.FetchUserByEmail(userData.Email)
	if err == nil {
		barf.Response(w).Status(http.StatusConflict).JSON(shared.ErrorResponse{
			StatusCode: 409,
			Errors: []map[string]string{
				{
					"email": fmt.Sprintf(
						"User with email address %s already exist on the app.", existingUserData.Email,
					),
				},
			},
			Message: "User already exist",
		})
		return
	}

	hashPwd, err := shared.HashPassword(userData.Password)
	if err != nil {
		barf.Logger().Error(err.Error())
		return
	}

	userData.Password = hashPwd
	newUserData, err2 := userService.CreateUser(userData)

	if err2 != nil {
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

	tokenString, err3 := shared.GenerateJWT(newUserData)
	if err3 != nil {
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
		Data: map[string]map[string]interface{}{
			"user": {
				"email":      userData.Email,
				"authToken":  tokenString,
				"_createdAt": newUserData.CreatedAt,
				"_updatedAt": newUserData.UpdatedAt,
			},
		},
	})
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(shared.UserSignupData{}).(*shared.UserSignupData)

	if userData == nil {
		barf.Logger().Error("no user data")
		barf.Response(w).Status(http.StatusBadRequest).JSON(shared.ErrorResponse{
			StatusCode: 400,
			Errors: []map[string]string{
				{
					"email": "Requires `email` field.",
				},
				{
					"password": "Requires `password` field.",
				},
			},
			Message: "Invalid request body.",
		})
		return
	}

	existingUserData, err := userService.FetchUserByEmail(userData.Email)

	if err != nil {
		barf.Response(w).Status(http.StatusNotFound).JSON(shared.ErrorResponse{
			StatusCode: 404,
			Errors: []map[string]string{
				{
					"email": fmt.Sprintf(
						"No User with email address %s on the app.", userData.Email,
					),
				},
			},
			Message: "No User found.",
		})
		return
	}

	if existingUserData != nil {
		if !shared.PasswordMatches(existingUserData.Password, userData.Password) {
			barf.Response(w).Status(http.StatusBadRequest).JSON(shared.ErrorResponse{
				StatusCode: 400,
				Errors: []map[string]string{
					{
						"message": "Login failed! credentials not correct.",
					},
				},
				Message: "Login failed! credentials not correct",
			})
			return
		}
	}

	tokenString, err := shared.GenerateJWT(existingUserData)
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
		Message:    "User Auth was successful.",
		Data: map[string]map[string]string{
			"user": {
				"email":     existingUserData.Email,
				"authToken": tokenString,
			},
		},
	})
}
