package middleware

import (
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"

)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
            }
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou não fornecido"})
            return
        }

        c.Next()
    }
}