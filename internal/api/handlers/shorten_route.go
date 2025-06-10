package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sameepkat/ushort/internal/service"
)

type ShortenRequest struct {
	URL       string    `json:"url" binding:"required"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ShortenResponse struct {
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	ExpiresAt   *time.Time `json:"exipres_at,omitempty"`
}

func Shorten(c *gin.Context, urlService *service.URLService) {
	var req ShortenRequest
	raw, exists := c.Get("sanitizedURL")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req, ok := raw.(ShortenRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request type"})
		return
	}

	url, err := urlService.CreateShortURL(c.Request.Context(), req.URL, nil, req.ExpiresAt)
	if err != nil {
		switch err {
		case service.ErrRateLimitExceeded:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		case service.ErrInvalidURL:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL provided"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		}

		return
	}

	shortURL := fmt.Sprintf("http://%s/%s", c.Request.Host, url.ShortCode)

	c.JSON(http.StatusOK, ShortenResponse{
		ShortURL:    shortURL,
		OriginalURL: url.OriginalURL,
		ExpiresAt:   &url.ExpiresAt,
	})

}
