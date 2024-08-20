package controllers

import (
	"net/http"
	"os"
	"server/config"
	"server/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
    EmailOrUsername string `json:"login"`
    Password        string `json:"password"`
}

func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    var user models.User
    if err := config.DB.Where("email = ? OR username = ?", req.EmailOrUsername, req.EmailOrUsername).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "ErrorFindLogin"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    token, err := generateJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    loginResponse := models.LoginResponse{
        User:  models.User{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
		},
        Token: token,
    }
    c.JSON(http.StatusOK, loginResponse)
}

func generateJWT(user models.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    jwtSecret := os.Getenv("JWT_SECRET")
    return token.SignedString([]byte(jwtSecret))
}