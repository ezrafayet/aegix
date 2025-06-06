package repository

import (
	"aegix/internal/domain"
	"aegix/internal/providers"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

var _ domain.RefreshTokenRepository = &RefreshTokenRepository{}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(refreshToken domain.RefreshToken) error {
	result := r.db.Model(&domain.RefreshToken{}).Create(&refreshToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByToken(token string) (domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	result := r.db.Model(&domain.RefreshToken{}).Where("token = ?", token).First(&refreshToken)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return domain.RefreshToken{}, result.Error
	}
	if result.Error == gorm.ErrRecordNotFound {
		return domain.RefreshToken{}, providers.ErrNoRefreshToken
	}
	return domain.RefreshToken{}, nil
}

func (r *RefreshTokenRepository) GetValidRefreshTokensByUserID(userID string) ([]domain.RefreshToken, error) {
	var refreshTokens []domain.RefreshToken
	result := r.db.Model(&domain.RefreshToken{}).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&refreshTokens)
	if result.Error != nil {
		return []domain.RefreshToken{}, result.Error
	}
	return refreshTokens, nil
}

func (r *RefreshTokenRepository) CleanExpiredTokens(userID string) error {
	result := r.db.Model(&domain.RefreshToken{}).Where("user_id = ? AND expires_at < ?", userID, time.Now()).Delete(&domain.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteRefreshToken(token string) error {
	result := r.db.Model(&domain.RefreshToken{}).Where("token = ?", token).Delete(&domain.RefreshToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
