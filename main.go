package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// "webserver/controller"
	// "webserver/database"
	// "webserver/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func logRequestMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        url := c.Request.URL.Scheme + "://" + c.Request.Host + c.Request.RequestURI
        fmt.Printf("url: %s\n", url)
        fmt.Printf("metodo: %s\n", c.Request.Method)
        c.Next()
    }
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Erro carregando .env")
    }

    // connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
    // os.Getenv("DB_USER"),
    // os.Getenv("DB_PASSWORD"),
    // os.Getenv("DB_HOST"),
    // os.Getenv("DB_PORT"),
    // os.Getenv("DB_NAME"))

    // var err error
    // db, err = sql.Open("postgres", connStr)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer db.Close()

    // err = db.Ping()
    // if err != nil {
    //     log.Fatal("Erro ao conectar ao banco de dados:", err)
    // }

    port := os.Getenv("PORT")
    if port == "" {
        port = "4000"
    }

    // database.CreateTables(db)

    r := gin.Default()
    r.Use(logRequestMiddleware())

    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Servidor Iniciado")
    })

    // r.POST("/register", func(c *gin.Context) {
    //     controller.RegisterHandler(c.Writer, c.Request, db)
    // })

    // r.POST("/login", func(c *gin.Context) {
    //     controller.LoginHandler(c.Writer, c.Request, db)
    // })

    // authGroup := r.Group("/")
    // authGroup.Use(middleware.Authenticate())
    // {
    //     authGroup.PUT("/update/:id", func(c *gin.Context) {
    //         controller.UpdateHandler(c.Writer, c.Request, db)
    //     })

    //     authGroup.DELETE("/delete/:id", func(c *gin.Context) {
    //         controller.DeleteHandler(c.Writer, c.Request, db)
    //     })
    // }

    fmt.Printf("Servidor iniciado na porta %s...\n", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}