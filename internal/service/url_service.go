package service

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/sameepkat/ushort/internal/cache"
	"github.com/sameepkat/ushort/internal/encoding"
	"github.com/sameepkat/ushort/internal/models"
	"github.com/sameepkat/ushort/internal/ratelimit"
	"gorm.io/gorm"
)

var (
	ErrURLNotFound       = errors.New("URL not found")
	ErrInvalidURL        = errors.New("invalid URL")
	ErrRateLimitExceeded = errors.New("rate limit exceeded. slow down bruh")
)

type URLService struct {
	db          *gorm.DB
	cache       *cache.Cache
	rateLimiter *ratelimit.RateLimiter
	mu          sync.Mutex
}

func NewURLService(db *gorm.DB, redisURL string) (*URLService, error) {
	cache, err := cache.NewCache(redisURL)
	if err != nil {
		return nil, err
	}

	rateLimiter := ratelimit.NewRateLimiter(1, time.Second)

	return &URLService{
		db:          db,
		cache:       cache,
		rateLimiter: rateLimiter,
	}, nil
}

func (s *URLService) CreateShortURL(ctx context.Context, originalURL string, userID *uint64, customCode string, expiresAt time.Time) (*models.URL, error) {
	if !s.rateLimiter.Allow(userID) {
		return nil, ErrRateLimitExceeded
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if customCode != "" {
		var existingURL models.URL
		if err := s.db.Where("short_code = ? ", customCode).First(&existingURL).Error; err == nil {
			return nil, errors.New("custom code already in use")
		}
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	url := &models.URL{
		OriginalURL: originalURL,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	if !expiresAt.IsZero() {
		url.ExpiresAt = expiresAt
	}

	if err := tx.Create(url).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if customCode != "" {
		url.ShortCode = customCode
	} else {
		url.ShortCode = encoding.Encode(url.ID)
	}

	if err := tx.Save(url).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, url.ShortCode, url, 1*time.Hour); err != nil {
		log.Printf("Failed to cache new URL: %v", err)
	}

	return url, nil
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (*models.URL, error) {
	if url, err := s.cache.Get(ctx, shortCode); err != nil {
		return nil, err
	} else if url != nil {
		return url, nil
	}
	var url models.URL

	if err := s.db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrURLNotFound
		}
		return nil, err
	}

	if !url.ExpiresAt.IsZero() && url.ExpiresAt.Before(time.Now()) {
		return nil, ErrURLNotFound
	}

	s.db.Model(&url).UpdateColumn("click_count", gorm.Expr("click_count + ?", 1))

	if err := s.cache.Set(ctx, shortCode, &url, 1*time.Hour); err != nil {
		log.Printf("Failed to cache URL: %v", err)
	}

	return &url, nil
}

func (s *URLService) GetUserURLs(userID uint64) ([]models.URL, error) {
	var urls []models.URL
	if err := s.db.Where("user_id = ?", userID).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}

func (s *URLService) DeleteURL(ctx context.Context, id uint64, userID uint64) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.URL{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrURLNotFound
	}

	shortCode := encoding.Encode(id)
	if err := s.cache.Delete(ctx, shortCode); err != nil {
		log.Printf("Failed to delete URL from cache: %v", err)
	}

	return nil
}

func (s *URLService) Close() error {
	return s.cache.Close()
}
