package postgres_gorm

import (
	"Docker/entity"
	"Docker/service"
	"context"

	"gorm.io/gorm"
)

// urlRepository is the implementation of IUrlRepository
type urlRepository struct {
	db *gorm.DB
}

// NewUrlRepository creates a new instance of urlRepository
func NewUrlRepository(db *gorm.DB) service.IUrlRepository {
	return &urlRepository{db: db}
}

// CreateUrl creates a new URL record in the database
func (r *urlRepository) CreateUrl(ctx context.Context, url *entity.Url) (entity.Url, error) {
	if err := r.db.WithContext(ctx).Create(url).Error; err != nil {
		return entity.Url{}, err
	}
	return *url, nil
}

// GetUrlByShortUrl retrieves a URL record by its short URL
func (r *urlRepository) GetUrlByShortUrl(ctx context.Context, shortUrl string) (entity.Url, error) {
	var url entity.Url
	if err := r.db.Where("short_url = ?", shortUrl).First(&url).Error; err != nil {
		return entity.Url{}, err
	}
	return url, nil
}
