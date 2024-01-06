package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ProductAPI/models"
	"strings"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context)  {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Fill all fields please!"})
		return
	}

	if !strings.Contains(user.Email, "@") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Fill the email corectly please!"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	if models.DB.Model(&user).Where("email = ?", user.Email).First(&user).RowsAffected > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email not available!"})
		return
	}
	
	models.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Login(c *gin.Context)  {
	var user models.User
	var userInput models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error})
		return
	}

	if models.DB.Model(&user).Where("email = ?", user.Email).First(&user).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Register first please!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInput.Password), []byte(user.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password wrong!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login success"})
}