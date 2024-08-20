package client

import (
	"net/http"
	"server/config"
	"server/models"

	"github.com/gin-gonic/gin"
)

func AddClient(c *gin.Context){
	var req models.Client

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}

	userID := c.MustGet("userID").(int)

	var user models.User
    if err := config.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
        return
    }

	client := models.Client{
		Name:  req.Name,
		Document: req.Document,
		Email: req.Email,
		Phone: models.Phone{
			CountryCode: req.Phone.CountryCode,
			AreaCode: req.Phone.AreaCode,
			Number: req.Phone.Number,
		},
		Address: models.Address{
			Street: req.Address.Street,
			Number: req.Address.Number,
			Complement: req.Address.Complement,
			District: req.Address.District,
			City: req.Address.City,
			State: req.Address.State,
			Country: req.Address.Country,
			Zipcode:req.Address.Zipcode,
		},
		UserID: userID,
	}

	if err := config.DB.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
	}

	if err := config.DB.Preload("User").Preload("Phone").Preload("Address").First(&client, client.ID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load client associations"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"Message": "clienteCreate"})
}