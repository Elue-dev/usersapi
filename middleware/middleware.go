package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/elue-dev/usersapi/utils"
)

type contextKey string

const userKey contextKey = "user"


func VerifyToken(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authorizationHeader := r.Header.Get("Authorization")

        if authorizationHeader == "" {
            utils.RespondWithError(w, http.StatusUnauthorized, "Missing Authorization Header")
            return
        }

        tokenString := strings.Split(authorizationHeader, " ")[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil
        })

        if err != nil {
            utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Token")
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Store the user information in the request context
            ctx := context.WithValue(r.Context(), userKey, claims["user"])
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Token")
        }
    }
}

func VerifyTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return VerifyToken(next)
}

