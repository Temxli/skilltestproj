package middleware

import (
	"auth-service/internal/database"
	"auth-service/internal/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth middleware
func RequireAuth(c *gin.Context) {
	// ensure DB is ready
	if database.DB == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration: database not initialized"})
		return
	}

	// Get the cookie from the request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No token provided"})
		return
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server misconfiguration"})
		return
	}

	// Parse with claims to get typed access to fields
	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token"})
		return
	}

	// Check expiration (allow numeric types safely)
	if expRaw, ok := claims["exp"]; ok {
		var expUnix int64
		switch v := expRaw.(type) {
		case float64:
			expUnix = int64(v)
		case int64:
			expUnix = v
		// case jsonNumberString:
		// 	// unreachable here but left for clarity
		case string:
			if parsed, perr := strconv.ParseInt(v, 10, 64); perr == nil {
				expUnix = parsed
			}
		}
		if expUnix != 0 && time.Now().Unix() > expUnix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Token expired"})
			return
		}
	}

	// Convert "sub" to uint primary key value
	var subID uint
	switch v := claims["sub"].(type) {
	case float64:
		subID = uint(v)
	case int64:
		subID = uint(v)
	case string:
		if parsed, perr := strconv.ParseUint(v, 10, 64); perr == nil {
			subID = uint(parsed)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - invalid subject claim"})
			return
		}
	default:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - invalid subject claim"})
		return
	}

	// Find user in database
	var user models.User
	if err := database.DB.First(&user, subID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - User not found"})
		return
	}

	// Attach user to request context
	c.Set("user", user)

	// Continue to the next handler
	c.Next()
}
