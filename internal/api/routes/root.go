package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sameepkat/ushort/internal/api/handlers"
	"github.com/sameepkat/ushort/internal/api/middleware"
	"github.com/sameepkat/ushort/internal/service"
)

func SetupRoutes(c *gin.RouterGroup, urlService *service.URLService, userService *service.UserService) {
	c.POST("/shorten", middleware.UrlSanitizer, func(c *gin.Context) { handlers.Shorten(c, urlService) })
	c.GET("/:short_url", func(c *gin.Context) { handlers.GetURL(c, urlService) })
	c.POST("/login", handlers.LoginHandler(userService))
	c.POST("/register", handlers.SignupHandler(userService))
}
