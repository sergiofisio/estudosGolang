package main

import (
	"encoding/json"
	"estudargolang/controllers"
	"estudargolang/middleware"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Rota chamada: %s, Método: %s", r.URL.Path, r.Method)
        next.ServeHTTP(w, r)
    })
}

func CustomHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    log.Printf("Rota inválida chamada: %s, Método: %s", r.URL.Path, r.Method)
    response := map[string]bool{"route": false}
    json.NewEncoder(w).Encode(response)
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
    }

    port, host := os.Getenv("PORT"), os.Getenv("HOST")
    if port == "" {
        log.Fatal("Porta não definida no arquivo .env")
    }
    if host == "" {
        host = "localhost"
    }

    r := mux.NewRouter()

    r.Use(LoggingMiddleware)

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Rota chamada: %s", r.URL.Path)
        fmt.Fprintln(w, "Servidor iniciado e rodando!")
    }).Methods("GET")

    r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Rota chamada: %s", r.URL.Path)
        controllers.Register(w, r)
    }).Methods("POST")

    r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Rota chamada: %s", r.URL.Path)
        controllers.Login(w, r)
    }).Methods("POST")

    r.HandleFunc("/update/{id}", middleware.ValidateTokenMiddleware(controllers.Update)).Methods("PUT")
    r.HandleFunc("/delete/{id}", middleware.ValidateTokenMiddleware(controllers.Delete)).Methods("DELETE")

    r.NotFoundHandler = http.HandlerFunc(CustomHandler)
    r.MethodNotAllowedHandler = http.HandlerFunc(CustomHandler)

    log.Printf("Server started on http://%s:%s", host, port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}