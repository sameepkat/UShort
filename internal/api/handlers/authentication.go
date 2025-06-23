package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sameepkat/ushort/internal/models"
	"github.com/sameepkat/ushort/internal/service"
	"github.com/sameepkat/ushort/internal/utils"
)

var jwtKey = []byte("asupersecretkey")

func LoginHandler(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
			return
		}

		user, err := userService.Authenticate(c, input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := models.Claims{
			UserID: user.ID,
			Role:   user.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Subject:   user.Email,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}

		c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "", false, false)
		c.JSON(http.StatusOK, gin.H{"message": "user logged in"})
	}
}

func SignupHandler(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := userService.CreateUser(c, input.Email, input.Password)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user created"})
	}
}

func Continue(c *gin.Context) {
	cookie, err := c.Cookie("token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := utils.ParseToken(cookie)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "continue", "role": claims.Role})

}
