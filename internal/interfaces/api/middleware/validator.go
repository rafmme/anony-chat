package middleware

import (
	"net/http"

	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/pkg/shared"
)

var (
	UserData *shared.UserSignupData
)

func ValidateSignupData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserData = new(shared.UserSignupData)
		err := barf.Request(r).Body().Format(UserData)

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

		validationResult := shared.SignUpDataValidator(*UserData)

		if len(validationResult) > 0 {
			barf.Response(w).Status(http.StatusBadRequest).JSON(shared.ErrorResponse{
				StatusCode: 400,
				Errors:     validationResult,
				Message:    "Invalid user signup data.",
			})
			return
		}

		next.ServeHTTP(w, r)
	}
}
