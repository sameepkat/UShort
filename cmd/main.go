package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sameepkat/ushort/internal/api/routes"
	"github.com/sameepkat/ushort/internal/database"
	"github.com/sameepkat/ushort/internal/service"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	config := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "ushort"),
		Password: getEnv("DB_PASSWORD", "randompassword"),
		DBName:   getEnv("DB_NAME", "ushort"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
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
