package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sameepkat/ushort/internal/api/routes"
	"github.com/sameepkat/ushort/internal/database"
	"github.com/sameepkat/ushort/internal/service"
)

func main() {
	// Enable this in release mode
	gin.SetMode(gin.ReleaseMode)

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Database URL not set")
	}

	u, err := url.Parse(databaseURL)
	if err != nil {
		log.Fatalf("Invalid DATABASE_URL: %v", err)
	}
	userInfo := u.User
	password, _ := userInfo.Password()

	config := database.Config{
		Host:     u.Hostname(),
		Port:     u.Port(),
		User:     userInfo.Username(),
		Password: password,
		DBName:   strings.TrimPrefix(u.Path, "/"),
		SSLMode:  u.Query().Get("sslmode"),
	}
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379/0")
	db, err := database.NewDB(config)
	if err != nil {
		log.Fatal(err)
	}
	urlService, err := service.NewURLService(db, redisURL)
	if err != nil {
		log.Fatal(err)
	}
	defer urlService.Close()

	userService := service.NewUserService(db)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://ushortlink.netlify.app/"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	root := router.Group("/")
	{
		routes.SetupRoutes(root, urlService, userService)
	}

	router.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
