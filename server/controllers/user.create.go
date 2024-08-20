package controllers

import (
	"net/http"
	"server/config"
	"server/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

func CreateUser(c *gin.Context) {
    var input struct {
        Name     string `json:"name" validate:"required"`
        Email    string `json:"email" validate:"required,email"`
        Username string `json:"username" validate:"required"`
        Password string `json:"password" validate:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := validate.Struct(input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var existingUser models.User
    if err := config.DB.Where("email = ? OR username = ?", input.Email, input.Username).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "emailUsername"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "errorPassword"})
        return
    }

    user := models.User{
        Name:     input.Name,
        Email:    input.Email,
        Username: input.Username,
        Password: string(hashedPassword),
    }
    config.DB.Create(&user)

    c.JSON(http.StatusOK, gin.H{"message": "userCreate"})
}