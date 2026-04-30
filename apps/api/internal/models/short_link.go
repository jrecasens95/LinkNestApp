package models

import "time"

type ShortLink struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Code        string       `gorm:"size:6;uniqueIndex;not null" json:"code"`
	OriginalURL string       `gorm:"not null" json:"original_url"`
	Title       *string      `gorm:"size:255" json:"title,omitempty"`
	ClicksCount uint         `gorm:"not null;default:0" json:"clicks_count"`
	IsActive    bool         `gorm:"not null;default:true" json:"is_active"`
	ExpiresAt   *time.Time   `json:"expires_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	ClickEvents []ClickEvent `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

type ClickEvent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ShortLinkID uint      `gorm:"not null;index" json:"short_link_id"`
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	Referer     string    `gorm:"type:text" json:"referer"`
	IPAddress   string    `gorm:"size:64" json:"ip_address"`
	CreatedAt   time.Time `json:"created_at"`
}
