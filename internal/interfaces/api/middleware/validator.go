package middleware

import (
	"fmt"
	"net/http"

	"github.com/opensaucerer/barf"
	infrastructure "github.com/rafmme/anony-chat/internal/infrastructure/persistence"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func ValidateSignupData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userSignupData := new(shared.UserSignupData)
		err := barf.Request(r).Body().Format(userSignupData)

		if err != nil {
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

		validationResult := shared.SignUpDataValidator(*userSignupData)

		if len(validationResult) > 0 {
			barf.Response(w).Status(http.StatusBadRequest).JSON(shared.ErrorResponse{
				StatusCode: 400,
				Errors:     validationResult,
				Message:    "Invalid user signup data.",
			})
			return
		}

		existingUserData := infrastructure.UserRepo.FindByEmail(userSignupData.Email)

		if len(existingUserData.Email) > 0 {
			barf.Response(w).Status(http.StatusConflict).JSON(shared.ErrorResponse{
				StatusCode: 409,
				Errors: []map[string]string{
					{
						"email": fmt.Sprintf(
							"User with email address %s already exist on the app.", userSignupData.Email,
						),
					},
				},
				Message: "User already exist",
			})
			return
		}

		next(w, r)
	}
}
