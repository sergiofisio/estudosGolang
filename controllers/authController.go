package controllers

import (
	"database/sql"
	"encoding/json"
	"estudargolang/config"
	"estudargolang/models"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var jwtKey = []byte("JWT_SECRET_KEY")
func init() {
    db = config.Connect()
}

func Register(w http.ResponseWriter, r *http.Request) {

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var exists bool
    err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
    if err != nil {
        http.Error(w, "Failed to query database", http.StatusInternalServerError)
        return
    }
    if exists {
        http.Error(w, "User already exists", http.StatusConflict) // HTTP 409 Conflict
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
        return
    }

    _, err = db.Exec("INSERT INTO users (name, email, document, username, password) VALUES ($1, $2, $3, $4, $5)", user.Name, user.Email, user.Document, user.Username, string(hashedPassword))
    if err != nil {
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintln(w, "User registered successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {

    var loginDetails struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    err := json.NewDecoder(r.Body).Decode(&loginDetails)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var user models.User

    err = db.QueryRow("SELECT id, name, email, document, username FROM users WHERE username = $1", loginDetails.Username).Scan(&user.ID, &user.Name, &user.Email, &user.Document, &user.Username)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to query user", http.StatusInternalServerError)
        }
        return
    }

    var hashedPassword string
    err = db.QueryRow("SELECT password FROM users WHERE username = $1", loginDetails.Username).Scan(&hashedPassword)
    if err != nil {
        http.Error(w, "Failed to query user password", http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginDetails.Password))
    if err != nil {
        http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &jwt.StandardClaims{
        Subject:   fmt.Sprint(user.ID),
        ExpiresAt: expirationTime.Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    response := struct {
        User  models.User `json:"user"`
        Token string      `json:"token"`
    }{
        User:  user,
        Token: tokenString,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func Update(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)
    id := vars["id"]

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.Exec("UPDATE users SET name = $1, email = $2, document = $3, username = $4, password = $5 WHERE id = $6", user.Name, user.Email, user.Document, user.Username, user.Password, id)
    if err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "User updated successfully")
}

func Delete(w http.ResponseWriter, r *http.Request) {

    vars := mux.Vars(r)
    id := vars["id"]

    _, err := db.Exec("DELETE FROM users WHERE id = $1", id)
    if err != nil {
        http.Error(w, "Failed to delete user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "User deleted successfully")
}