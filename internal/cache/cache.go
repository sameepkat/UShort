package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sameepkat/ushort/internal/models"
)

type Cache struct {
	client *redis.Client
}

type cacheItem struct {
	URL       *models.URL `json:"url"`
	ExpiresAt time.Time   `json:"expires_at"`
}

func NewCache(redisURL string) (*Cache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{
		client: client,
	}, nil
}

func (c *Cache) Get(ctx context.Context, shortCode string) (*models.URL, error) {
	data, err := c.client.Get(ctx, shortCode).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var item cacheItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, err
	}

	if time.Now().After(item.ExpiresAt) {
		c.client.Del(ctx, shortCode)
		return nil, nil
	}

	return item.URL, nil
}

func (c *Cache) Set(ctx context.Context, shortCode string, url *models.URL, ttl time.Duration) error {
	item := cacheItem{
		URL:       url,
		ExpiresAt: time.Now().Add(ttl),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return nil
	}

	return c.client.Set(ctx, shortCode, data, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, shortCode string) error {
	return c.client.Del(ctx, shortCode).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}

// func (c *Cache) GetStats(ctx context.Context) (map[string]interface{}, error) {
// 	info, err := c.client.Info(ctx).Result()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return info, nil
// }
