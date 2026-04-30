package models

import "time"

type ShortLink struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Code        string     `gorm:"size:6;uniqueIndex;not null" json:"code"`
	OriginalURL string     `gorm:"not null" json:"original_url"`
	Title       *string    `gorm:"size:255" json:"title,omitempty"`
	ClicksCount uint       `gorm:"not null;default:0" json:"clicks_count"`
	IsActive    bool       `gorm:"not null;default:true" json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
