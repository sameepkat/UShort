package models

import (
	"time"
)

type URL struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	ShortCode   string    `json:"short_code" gorm:"uniqueIndex;not null"`
	OriginalURL string    `json:"original_url" gorm:"not null"`
	UserID      *uint64   `json:"user_id" gorm:"index"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	ExpiresAt   time.Time `json:"uint64" gorm:"default:0"`
}
