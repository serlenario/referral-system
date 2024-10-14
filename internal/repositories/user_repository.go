// repositories/user_repository.go
package repositories

import (
	"github.com/serlenario/referral-system/internal/models"
	"gorm.io/gorm"
)

// UserRepository определяет методы для взаимодействия с пользователями в базе данных
type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByReferralCode(code string) (*models.User, error)
	Update(user *models.User) error
}

// userRepo реализует интерфейс UserRepository
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository создаёт новый экземпляр userRepo
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

// Create добавляет нового пользователя в базу данных
func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByEmail получает пользователя по email
func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID получает пользователя по ID
func (r *userRepo) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByReferralCode получает пользователя по реферальному коду
func (r *userRepo) GetByReferralCode(code string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("referral_code = ?", code).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update сохраняет изменения пользователя в базе данных
func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}
