package handlers

import (
	"fmt"
	"net/http"

	"github.com/opensaucerer/barf"
	infrastructure "github.com/rafmme/anony-chat/internal/infrastructure/persistence"
	"github.com/rafmme/anony-chat/internal/interfaces/api/middleware"
	domain "github.com/rafmme/anony-chat/pkg/domain/user"
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

	hashPwd, err := shared.HashPassword(middleware.UserData.Password)
	if err != nil {
		barf.Logger().Error(err.Error())
		return
	}

	middleware.UserData.Password = hashPwd
	userData := infrastructure.UserRepo.Save(middleware.UserData)
	tokenString, err := shared.GenerateJWT(userData)

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

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	userData := new(shared.UserSignupData)
	err := barf.Request(r).Body().Format(userData)

	if err != nil {
		barf.Logger().Error(err.Error())
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

	existingUserData := infrastructure.UserRepo.Find(
		&domain.User{
			Email: userData.Email,
		},
	)

	if len(existingUserData.Email) > 1 {
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

	if len(existingUserData.Email) < 1 {
		barf.Response(w).Status(http.StatusNotFound).JSON(shared.ErrorResponse{
			StatusCode: 404,
			Errors: []map[string]string{
				{
					"email": fmt.Sprintf(
						"No User with email address %s on the app.", middleware.UserData.Email,
					),
				},
			},
			Message: "No User found.",
		})
		return
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
