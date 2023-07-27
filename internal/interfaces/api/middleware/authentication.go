package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			barf.Response(w).Status(http.StatusUnauthorized).JSON(shared.ErrorResponse{
				StatusCode: 401,
				Errors: []map[string]string{
					{
						"auth": "Unauthorized.",
					},
				},
				Message: "Unauthorized.",
			})
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		env, err := shared.LoadEnvVars()

		if err != nil {
			barf.Logger().Error(err.Error())
			os.Exit(1)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(env.SecretKey), nil
		})

		if err != nil || !token.Valid {
			barf.Response(w).Status(http.StatusUnauthorized).JSON(shared.ErrorResponse{
				StatusCode: 401,
				Errors: []map[string]string{
					{
						"auth": "Unauthorized.",
					},
				},
				Message: "Unauthorized.",
			})
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			barf.Response(w).Status(http.StatusUnauthorized).JSON(shared.ErrorResponse{
				StatusCode: 401,
				Errors: []map[string]string{
					{
						"auth": "Unauthorized.",
					},
				},
				Message: "Unauthorized.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
