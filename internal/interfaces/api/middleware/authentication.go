package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/opensaucerer/barf"
	"github.com/rafmme/anony-chat/pkg/shared"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := shared.GetAuthToken(r, "ws")

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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
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

		claims, ok := token.Claims.(jwt.MapClaims)
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

		ctx := context.WithValue(
			r.Context(), shared.AuthData{},
			claims,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
