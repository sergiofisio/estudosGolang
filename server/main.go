package main

import (
	"fmt"
	"net/http"
	"server/config"
	"server/controllers"
	"server/controllers/client"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"server_started": true,
			"message":        "Servidor iniciado",
		})
	})

	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users/:id", controllers.FindUser)
		protected.POST("/addclient", client.AddClient)

		protected.GET("/clients/:id", client.FindClient)
	}

	return r
}

func main() {
	config.LoadConfig()
	config.ConnectDatabase()

	r := setupRouter()
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	r.Run(port)
}
