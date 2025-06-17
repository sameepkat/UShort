package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sameepkat/ushort/internal/service"
)

type URLResonse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func GetURL(urlService *service.URLService) gin.HandlerFunc {
	return func(c *gin.Context) {

		short_code := c.Param("short_url")

		url, err := urlService.GetOriginalURL(c, short_code)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, URLResonse{
			ShortURL:    url.ShortCode,
			OriginalURL: url.OriginalURL,
		})
	}
}
