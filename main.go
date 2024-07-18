package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux" // Importe o gorilla/mux
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sergiofisio/estudosGolang/controller"
	"github.com/sergiofisio/estudosGolang/middleware"
)

var db *sql.DB

func logRequestMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        url := r.URL.Scheme + "://" + r.Host + r.RequestURI
        fmt.Printf("url: %s\n", url)
        fmt.Printf("metodo: %s\n", r.Method)
        next.ServeHTTP(w, r)
    })
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Erro carregando .env")
    }

    connStr := "user=postgres password=123456 dbname=golang host=localhost port=5434 sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal("Erro ao conectar ao banco de dados:", err)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "4000"
    }

    CreateTables(db)

    r := mux.NewRouter()

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Servidor Iniciado")
    }).Methods("GET")

    r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        controller.RegisterHandler(w, r, db)
    }).Methods("POST")

    r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        controller.LoginHandler(w, r, db)
    }).Methods("POST")

    r.Handle("/update/{id}", middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        controller.UpdateHandler(w, r, db)
    }))).Methods("PUT")

    r.Handle("/delete/{id}", middleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        controller.DeleteHandler(w, r, db)
    }))).Methods("DELETE")

    wrappedMux := logRequestMiddleware(r) 

    fmt.Printf("Servidor iniciado na porta %s...\n", port)
    if err := http.ListenAndServe(":"+port, wrappedMux); err != nil {
        log.Fatal(err)
    }
}