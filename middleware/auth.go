package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

)

var jwtKey = []byte("JWT_SECRET_KEY")

func ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        claims := &jwt.StandardClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        next(w, r)
    }
}