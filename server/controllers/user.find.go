package controllers

import (
	"net/http"
	"server/config"
	"server/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FindUser(c *gin.Context) {
    id := c.Param("id")
    userID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var user models.User
    if err := config.DB.Preload("Clients").Preload("Clients.Phone").Preload("Clients.Address").First(&user, userID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }

    for i := range user.Clients {
        user.Clients[i].UserID = 0
        user.Clients[i].User = models.User{}
    }
    

    c.JSON(http.StatusOK, models.User{
        ID:       user.ID,
        Name:     user.Name,
        Email:    user.Email,
        Username: user.Username,
        Clients:  user.Clients, 
    })
}