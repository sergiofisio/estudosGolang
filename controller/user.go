package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
    "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Address struct {
	Street       string `json:"street"`
	Number       int    `json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	ZipCode      string `json:"zip_code"`
}

type Phone struct {
	PhoneType   string `json:"phone_type"`
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
	PhoneNumber string `json:"phone_number"`
}

type Purchase struct {
	ProductName  string  `json:"product_name"`
	Amount       float64 `json:"amount"`
	PurchaseDate string  `json:"purchase_date"`
}

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

func setupLogger() {
    logFile, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("Erro ao abrir o arquivo de log:", err)
    }
    log.SetOutput(logFile)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
    setupLogger()
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        LogError(w, "RegisterHandler", "Erro ao decodificar o corpo da requisição", err, http.StatusBadRequest)
        return
    }

    if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
        LogError(w, "RegisterHandler", "Erro ao encriptar a senha", err, http.StatusBadRequest)
    } else {
        user.Password = string(hashedPassword)
        if err := SaveUser(db, user.Name, user.Email, user.Document, user.Username, user.Password); err != nil {
            LogError(w, "RegisterHandler", "Erro ao salvar o usuário no banco de dados", err, http.StatusBadRequest)
            return
        }

        user.Password = ""
        SendJSONResponse(w, http.StatusCreated, struct {
            User    User `json:"user"`
            Message string      `json:"message"`
        }{
            User:    user,
            Message: "Usuário registrado com sucesso",
        })
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var loginDetails struct {
        UsernameOrEmail string `json:"usernameOrEmail"`
        Password        string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
        LogError(w, "LoginHandler", "Erro ao decodificar o corpo da requisição", err, http.StatusBadRequest)
        return
    }

    var user User
    if err := db.QueryRow(`SELECT id, name, email, document, username, password FROM users WHERE username = $1 OR email = $1`, loginDetails.UsernameOrEmail).Scan(&user.ID, &user.Name, &user.Email, &user.Document, &user.Username, &user.Password); err != nil {
        LogError(w, "LoginHandler", "Usuário não encontrado", err, http.StatusBadRequest)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password)); err != nil {
        LogError(w, "LoginHandler", "Senha inválida", err, http.StatusBadRequest)
        return
    }

    tokenString, err := GenerateJWTToken(user.Email)
    if err != nil {
        LogError(w, "LoginHandler", "Erro ao gerar o token", err, http.StatusBadRequest)
        return
    }

    user.Password = ""
    SendJSONResponse(w, http.StatusOK, struct {
        User  User `json:"user"`
        Token string      `json:"token"`
    }{
        User:  user,
        Token: tokenString,
    })
}

func UpdateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    vars := mux.Vars(r)
    userID := vars["id"]

    print(userID)

    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        LogError(w, "UpdateHandler", "Erro ao decodificar o corpo da requisição", err, http.StatusBadRequest)
        return
    }

    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
    if err != nil {
        LogError(w, "UpdateHandler", "Erro ao verificar a existência do usuário", err, http.StatusInternalServerError)
        return
    }
    if !exists {
        SendJSONResponse(w, http.StatusNotFound, "Usuário não encontrado")
        return
    }

    query := "UPDATE users SET "
    var updates []string
    var args []interface{}
    argID := 1

    if user.Name != "" {
        updates = append(updates, fmt.Sprintf("name = $%d", argID))
        args = append(args, user.Name)
        argID++
    }

    if user.Email != "" {
        updates = append(updates, fmt.Sprintf("email = $%d", argID))
        args = append(args, user.Email)
        argID++
    }

    if user.Document != "" {
        updates = append(updates, fmt.Sprintf("document = $%d", argID))
        args = append(args, user.Document)
        argID++
    }

    if user.Username != "" {
        updates = append(updates, fmt.Sprintf("username = $%d", argID))
        args = append(args, user.Username)
        argID++
    }

    if user.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            LogError(w, "UpdateHandler", "Erro ao encriptar a senha", err, http.StatusInternalServerError)
            return
        }
        updates = append(updates, fmt.Sprintf("password = $%d", argID))
        args = append(args, string(hashedPassword))
        argID++
    }

    if len(updates) == 0 {
        SendJSONResponse(w, http.StatusBadRequest, "Nenhum campo para atualizar")
        return
    }

    query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", argID)

    args = append(args, userID)

    _, err = db.Exec(query, args...)
    if err != nil {
        LogError(w, "UpdateHandler", "Erro ao atualizar o usuário no banco de dados", err, http.StatusInternalServerError)
        return
    }

    SendJSONResponse(w, http.StatusOK, "Usuário atualizado com sucesso")
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    userID := r.URL.Query().Get("id")
    if userID == "" {
        http.Error(w, "ID do usuário é obrigatório", http.StatusBadRequest)
        return
    }

    query := `DELETE FROM users WHERE id = $1;`
    result, err := db.Exec(query, userID)
    if err != nil {
        http.Error(w, "Erro ao deletar usuário", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, "Erro ao verificar usuário deletado", http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Nenhum usuário encontrado com o ID fornecido", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Usuário deletado com sucesso"))
}
