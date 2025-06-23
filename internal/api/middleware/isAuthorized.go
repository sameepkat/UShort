package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sameepkat/ushort/internal/utils"
)

func IsAuthorized(c *gin.Context) {
	cookie, err := c.Cookie("token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	claims, err := utils.ParseToken(cookie)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}

	c.Set("UserID", claims.UserID)

	c.Set("role", claims.Role)
	c.Next()
}
