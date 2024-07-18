package main

import (
	"estudargolang/controllers"
	"fmt"
	"log"
	"net/http"

)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rota chamada: %s", r.URL.Path)
        fmt.Fprintln(w, "Servidor iniciado e rodando!")
    })

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Rota chamada: %s", r.URL.Path)
        controllers.Register(w, r)
    })

    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Rota chamada: %s", r.URL.Path)
        controllers.Login(w, r)
    })

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}