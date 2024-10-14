package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Email          string         `gorm:"unique;not null" json:"email"`
	PasswordHash   string         `json:"-"`
	ReferralCode   string         `gorm:"unique" json:"referral_code"`
	ReferralExpiry time.Time      `json:"referral_expiry"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Referrals      []Referral     `json:"referrals,omitempty" gorm:"foreignKey:ReferredBy"`
}

type Referral struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ReferredID uint           `json:"referred_id"`
	ReferredBy uint           `json:"referred_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
