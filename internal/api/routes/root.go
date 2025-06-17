package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sameepkat/ushort/internal/api/handlers"
	"github.com/sameepkat/ushort/internal/api/middleware"
	"github.com/sameepkat/ushort/internal/service"
)

func SetupRoutes(c *gin.RouterGroup, urlService *service.URLService, userService *service.UserService) {
	c.POST("/shorten", middleware.IsAuthorized, middleware.UrlSanitizer, handlers.Shorten(urlService))
	c.GET("/:short_url", handlers.GetURL(urlService))
	c.POST("/login", handlers.LoginHandler(userService))
	c.POST("/register", handlers.SignupHandler(userService))
	c.GET("/favicon.ico", func(c *gin.Context) { c.Status(304) })
}
