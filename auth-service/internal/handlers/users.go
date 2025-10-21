package handlers

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, users)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user.Password == "" {
		c.JSON(400, gin.H{"error": "password required"})
		return
	}

	// hash the password before saving
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = string(hashed)

	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	user.Password = ""
	c.JSON(201, user)
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("name = ?", req.Name).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// build token
	secret := os.Getenv("SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration"})
		return
	}
	expiration := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": expiration.Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	// set cookie (HttpOnly). set Secure=true in production.
	secure := os.Getenv("ENV") == "production"
	maxAge := int(time.Until(expiration).Seconds())
	c.SetCookie("Authorization", signed, maxAge, "/", "", secure, true)

	// don't return password hash
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Logout(c *gin.Context) {
	// Clear the cookie by setting a past expiration
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
func ChangePassword(c *gin.Context) {

	id := c.Param("id")
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "User not found"})
		} else {
			c.JSON(500, gin.H{"error": result.Error.Error()})
		}
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(401, gin.H{"error": "Old password is incorrect"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash new password"})
		return
	}
	user.Password = string(hashed)
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Password changed successfully"})
}
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "User not found"})
		} else {
			c.JSON(500, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(200, user)
}
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "User not found"})
		} else {
			c.JSON(500, gin.H{"error": result.Error.Error()})
		}
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, user)
}
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(204)
}

func GetMyInfo(c *gin.Context) {
	userIface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	var user models.User
	switch u := userIface.(type) {
	case models.User:
		user = u
	case *models.User:
		if u == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user in context"})
			return
		}
		user = *u
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	// remove password before returning
	user.Password = ""
	c.JSON(http.StatusOK, user)
}
