package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message":"Servidor iniciado"})
	})

    router.POST("/register", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Registro bem-sucedido!"})
    })

    router.POST("/login", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Login bem-sucedido!"})
    })

    router.Run(":8080")
}