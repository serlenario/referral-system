package repositories

import (
	"github.com/serlenario/referral-system/internal/models"
	"gorm.io/gorm"
)

type ReferralRepository interface {
	Create(referral *models.Referral) error
	GetByReferrerID(referrerID uint) ([]models.Referral, error)
}

type referralRepo struct {
	db *gorm.DB
}

func NewReferralRepository(db *gorm.DB) ReferralRepository {
	return &referralRepo{db}
}

func (r *referralRepo) Create(referral *models.Referral) error {
	return r.db.Create(referral).Error
}

func (r *referralRepo) GetByReferrerID(referrerID uint) ([]models.Referral, error) {
	var referrals []models.Referral
	if err := r.db.Where("referred_by = ?", referrerID).Find(&referrals).Error; err != nil {
		return nil, err
	}
	return referrals, nil
}
