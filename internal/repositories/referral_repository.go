// repositories/referral_repository.go
package repositories

import (
	"github.com/serlenario/referral-system/internal/models"
	"gorm.io/gorm"
)

// ReferralRepository определяет методы для взаимодействия с рефералами в базе данных
type ReferralRepository interface {
	Create(referral *models.Referral) error
	GetByReferrerID(referrerID uint) ([]models.Referral, error)
}

// referralRepo реализует интерфейс ReferralRepository
type referralRepo struct {
	db *gorm.DB
}

// NewReferralRepository создаёт новый экземпляр referralRepo
func NewReferralRepository(db *gorm.DB) ReferralRepository {
	return &referralRepo{db}
}

// Create добавляет нового реферала в базу данных
func (r *referralRepo) Create(referral *models.Referral) error {
	return r.db.Create(referral).Error
}

// GetByReferrerID получает список рефералов по ID реферера
func (r *referralRepo) GetByReferrerID(referrerID uint) ([]models.Referral, error) {
	var referrals []models.Referral
	if err := r.db.Where("referred_by = ?", referrerID).Find(&referrals).Error; err != nil {
		return nil, err
	}
	return referrals, nil
}
