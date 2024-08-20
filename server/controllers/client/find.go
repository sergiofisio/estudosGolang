package client

import (
	"net/http"
	"server/config"
	"server/models"

	"github.com/gin-gonic/gin"
)

func FindClient(c *gin.Context){
	id := c.Param("id")

	var client models.Client

	if err := config.DB.Where("id = ?", id).First(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client not found"})
		return
	}

	if err := config.DB.Preload("User").Preload("Phone").Preload("Address").First(&client, client.ID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load client associations"})
        return
    }

	c.JSON(http.StatusOK, models.Client{
		ID: client.ID,
		Name: client.Name,
		Document: client.Document,
		Email: client.Email,
		PhoneID: client.PhoneID,
		Phone: client.Phone,
		AddressID: client.AddressID,
		Address: client.Address,
	})
}