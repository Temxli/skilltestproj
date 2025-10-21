package routes

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	authServiceURL    = "http://localhost:8082"
	productServiceURL = "http://localhost:8080"
	orderServiceURL   = "http://localhost:8081"
)

func ProxyAuthService(c *gin.Context) {
	forwardRequest(c, authServiceURL)
}

func ProxyProductService(c *gin.Context) {
	forwardRequest(c, productServiceURL)
}

func ProxyOrderService(c *gin.Context) {
	forwardRequest(c, orderServiceURL)
}

func forwardRequest(c *gin.Context, serviceURL string) {
	target := serviceURL + c.Request.URL.Path
	req, err := http.NewRequest(c.Request.Method, target, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header = c.Request.Header.Clone()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Request to service failed"})
		return
	}
	defer resp.Body.Close()

	// Копируем все заголовки кроме CORS
	for key, values := range resp.Header {
		if key == "Access-Control-Allow-Origin" ||
			key == "Access-Control-Allow-Credentials" ||
			key == "Access-Control-Allow-Headers" ||
			key == "Access-Control-Allow-Methods" {
			continue
		}
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	c.Writer.WriteHeader(resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	c.Writer.Write(body)
}
