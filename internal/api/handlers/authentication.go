package handlers

import (
	"log"
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
			Name:   user.Name,
			Email:  user.Email,
			Role:   user.Role,
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
		log.Println("cookie err")
		return
	}

	claims, err := utils.ParseToken(cookie)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		log.Println("parsetoken err")
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		log.Println("Claim.role = ", claims.Role)
		log.Printf("%s vs %s\n", claims.Role, "user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "user": gin.H{
		"name":  claims.Name,
		"email": claims.Email,
		"role":  claims.Role,
	},
	})

}
