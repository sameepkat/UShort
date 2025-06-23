package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sameepkat/ushort/internal/api/handlers"
)

func UrlSanitizer(c *gin.Context) {
	var req handlers.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body in middleware"})
		c.Abort()
		return
	}

	req.URL = strings.TrimSpace(req.URL)

	c.Set("sanitizedURL", req)
	c.Next()
}
