package middleware

import (
	"context"
	"net/http"

	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func ValidateSignupData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
					{
						"confirmPassword": "Requires `confirmPassword` field.",
					},
				},
				Message: "Invalid request body.",
			})
			return
		}

		validationResult := shared.AuthDataValidator("signup", *userData)

		if len(validationResult) > 0 {
			barf.Response(w).Status(http.StatusBadRequest).JSON(shared.ErrorResponse{
				StatusCode: 400,
				Errors:     validationResult,
				Message:    "Invalid user signup data.",
			})
			return
		}

		ctx := context.WithValue(r.Context(), shared.UserSignupData{}, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func ValidateAuthData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		validationResult := shared.AuthDataValidator("login", *userData)

		if len(validationResult) > 0 {
			barf.Response(w).Status(http.StatusBadRequest).JSON(shared.ErrorResponse{
				StatusCode: 400,
				Errors:     validationResult,
				Message:    "Invalid user login data.",
			})
			return
		}

		ctx := context.WithValue(r.Context(), shared.UserSignupData{}, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
