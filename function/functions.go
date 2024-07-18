package function

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func SaveUser(db *sql.DB, name, email, document, username, password string) error {
    query := `INSERT INTO users (name, email, document, username, password) VALUES ($1, $2, $3, $4, $5)`

    _, err := db.Exec(query, name, email, document, username, password)
    if err != nil {
        log.Printf("Erro ao inserir o usuário no banco de dados: %v", err)
        return err
    }

    return nil
}

func LogError(w http.ResponseWriter, functionName, message string, err error, statusCode int) {
    log.Printf("[%s] %s: %v\n", functionName, message, err)
    http.Error(w, message, statusCode)
}

func GenerateJWTToken(userEmail string) (string, error) {
    claims := &jwt.StandardClaims{
        Subject:   userEmail,
        ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)

    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
    w.WriteHeader(statusCode)
    w.Header().Set("Content-Type", "application/json")
    
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        LogError(w, "sendJSONResponse", "Erro ao codificar a resposta", err, http.StatusBadRequest)
    }
}