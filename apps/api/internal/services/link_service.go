package services

import (
	"errors"
	"strings"
	"time"

	"github.com/jrecasens95/link-nest/backend/internal/models"
	"github.com/jrecasens95/link-nest/backend/internal/utils"
	"gorm.io/gorm"
)

var (
	ErrLinkNotFound = errors.New("link not found")
	ErrLinkInactive = errors.New("link inactive")
	ErrLinkExpired  = errors.New("link expired")
)

type LinkService struct {
	db *gorm.DB
}

type UpdateLinkInput struct {
	Title    **string
	IsActive *bool
}

type ClickInput struct {
	UserAgent string
	Referer   string
	IPAddress string
}

type LinkStats struct {
	TotalClicks  uint
	RecentClicks []models.ClickEvent
	Referers     []RefererStat
}

type RefererStat struct {
	Referer string
	Count   int64
}

func NewLinkService(db *gorm.DB) *LinkService {
	return &LinkService{db: db}
}

func (s *LinkService) Create(originalURL string, title *string) (*models.ShortLink, error) {
	for range 10 {
		code, err := utils.GenerateCode(6)
		if err != nil {
			return nil, err
		}

		link := &models.ShortLink{
			Code:        code,
			OriginalURL: originalURL,
			Title:       title,
			IsActive:    true,
		}

		if err := s.db.Create(link).Error; err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				continue
			}
			return nil, err
		}

		return link, nil
	}

	return nil, errors.New("could not generate unique short code")
}

func (s *LinkService) List() ([]models.ShortLink, error) {
	var links []models.ShortLink
	if err := s.db.Order("created_at desc").Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

func (s *LinkService) GetByID(id uint) (*models.ShortLink, error) {
	var link models.ShortLink
	if err := s.db.First(&link, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}

	return &link, nil
}

func (s *LinkService) Update(id uint, input UpdateLinkInput) (*models.ShortLink, error) {
	link, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if len(updates) > 0 {
		if err := s.db.Model(link).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.GetByID(id)
}

func (s *LinkService) Delete(id uint) error {
	result := s.db.Delete(&models.ShortLink{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrLinkNotFound
	}

	return nil
}

func (s *LinkService) Stats(id uint) (*LinkStats, error) {
	link, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	var recentClicks []models.ClickEvent
	if err := s.db.
		Where("short_link_id = ?", id).
		Order("created_at desc").
		Limit(25).
		Find(&recentClicks).Error; err != nil {
		return nil, err
	}

	var referers []RefererStat
	if err := s.db.
		Model(&models.ClickEvent{}).
		Select("COALESCE(NULLIF(referer, ''), 'Direct') AS referer, COUNT(*) AS count").
		Where("short_link_id = ?", id).
		Group("COALESCE(NULLIF(referer, ''), 'Direct')").
		Order("count desc").
		Limit(10).
		Scan(&referers).Error; err != nil {
		return nil, err
	}

	return &LinkStats{
		TotalClicks:  link.ClicksCount,
		RecentClicks: recentClicks,
		Referers:     referers,
	}, nil
}

func (s *LinkService) Resolve(code string, click ClickInput) (*models.ShortLink, error) {
	var link models.ShortLink
	if err := s.db.Where("code = ?", code).First(&link).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}

	if !link.IsActive {
		return nil, ErrLinkInactive
	}

	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return nil, ErrLinkExpired
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		event := models.ClickEvent{
			ShortLinkID: link.ID,
			UserAgent:   click.UserAgent,
			Referer:     click.Referer,
			IPAddress:   click.IPAddress,
		}

		if err := tx.Create(&event).Error; err != nil {
			return err
		}

		return tx.Model(&link).UpdateColumn("clicks_count", gorm.Expr("clicks_count + ?", 1)).Error
	}); err != nil {
		return nil, err
	}

	return &link, nil
}
