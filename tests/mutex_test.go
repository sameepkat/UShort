package tests

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/sameepkat/ushort/internal/database"
	"github.com/sameepkat/ushort/internal/models"
	"github.com/sameepkat/ushort/internal/service"
)

func TestConcurrentURLCreation(t *testing.T) {
	redisURL := "redis://localhost:6379/0"

	config := database.Config{
		Host:     "localhost",
		Port:     "3307",
		User:     "admin",
		Password: "admin",
		DBName:   "urls",
		SSLMode:  "false",
	}
	db, err := database.NewDB(config)
	if err != nil {
		log.Fatal(err)
	}
	urlService, err := service.NewURLService(db, redisURL)
	if err != nil {
		log.Fatal(err)
	}
	defer urlService.Close()

	test_service, _ := service.NewURLService(db, redisURL)
	defer test_service.Close()

	var wg sync.WaitGroup
	urls := make(chan *models.URL, 10)
	errors := make(chan error, 10)

	// Create 10 URLs concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := test_service.CreateShortURL(context.Background(), "https://example.com", nil, time.Now().Add(1*time.Hour))
			if err != nil {
				errors <- err
			}
			urls <- url
		}()
	}

	wg.Wait()
	close(urls)
	close(errors)

	for url := range urls {
		if url.ShortCode == "" {
			t.Error("Empty Short Code generated")
		}
	}

	for err := range errors {
		t.Errorf("Error creating URL: %v", err)
	}
}
