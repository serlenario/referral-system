package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/serlenario/referral-system/internal/models"
	"github.com/serlenario/referral-system/internal/repositories"
	"github.com/serlenario/referral-system/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(email, password string) (*models.User, error)
	Authenticate(email, password string) (string, error)
	CreateReferralCode(userID uint, expiry time.Time) (*models.User, error)
	DeleteReferralCode(userID uint) (*models.User, error)
	GetReferralCodeByEmail(email string) (string, error)
	RegisterWithReferral(referralCode, email, password string) (*models.User, error)
	GetReferrals(userID uint) ([]models.Referral, error)
}

type userService struct {
	userRepo     repositories.UserRepository
	referralRepo repositories.ReferralRepository
	jwtSecret    string
}

func NewUserService(userRepo repositories.UserRepository, referralRepo repositories.ReferralRepository, jwtSecret string) UserService {
	return &userService{
		userRepo:     userRepo,
		referralRepo: referralRepo,
		jwtSecret:    jwtSecret,
	}
}

func (s *userService) Register(email, password string) (*models.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	return user, s.userRepo.Create(user)
}

func (s *userService) Authenticate(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) CreateReferralCode(userID uint, expiry time.Time) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	referralCode := uuid.New().String()

	user.ReferralCode = referralCode
	user.ReferralExpiry = expiry

	return user, s.userRepo.Update(user)
}

func (s *userService) DeleteReferralCode(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.ReferralCode = ""
	user.ReferralExpiry = time.Time{}

	return user, s.userRepo.Update(user)
}

func (s *userService) GetReferralCodeByEmail(email string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if !user.ReferralExpiry.IsZero() && user.ReferralExpiry.Before(time.Now()) {
		return "", errors.New("referral code expired")
	}

	if user.ReferralCode == "" {
		return "", errors.New("no referral code found")
	}

	return user.ReferralCode, nil
}

func (s *userService) RegisterWithReferral(referralCode, email, password string) (*models.User, error) {
	referrer, err := s.userRepo.GetByReferralCode(referralCode)
	if err != nil {
		return nil, errors.New("invalid referral code")
	}

	newUser, err := s.Register(email, password)
	if err != nil {
		return nil, err
	}

	referral := &models.Referral{
		ReferredID: newUser.ID,
		ReferredBy: referrer.ID,
	}

	if err := s.referralRepo.Create(referral); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) GetReferrals(userID uint) ([]models.Referral, error) {
	referrals, err := s.referralRepo.GetByReferrerID(userID)
	if err != nil {
		return nil, err
	}
	return referrals, nil
}
