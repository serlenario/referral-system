// internal/controllers/user_controller.go
package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/serlenario/referral-system/internal/models"
	"github.com/serlenario/referral-system/internal/services"
)

// UserController отвечает за обработку пользовательских запросов
type UserController struct {
	UserService services.UserService
}

// NewUserController создаёт новый экземпляр UserController
func NewUserController(userService services.UserService) *UserController {
	return &UserController{UserService: userService}
}

// RegisterRequest представляет структуру запроса для регистрации
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest представляет структуру запроса для аутентификации
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse представляет структуру ответа с JWT-токеном
type TokenResponse struct {
	Token string `json:"token"`
}

// CreateReferralRequest представляет структуру запроса для создания реферального кода
type CreateReferralRequest struct {
	Expiry time.Time `json:"expiry" binding:"required"`
}

// ReferralResponse представляет структуру ответа с реферальным кодом
type ReferralResponse struct {
	ReferralCode string    `json:"referral_code"`
	Expiry       time.Time `json:"expiry"`
}

// RegisterWithReferralRequest представляет структуру запроса для регистрации с реферальным кодом
type RegisterWithReferralRequest struct {
	ReferralCode string `json:"referral_code" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
}

// ReferralsResponse представляет структуру ответа со списком рефералов
type ReferralsResponse struct {
	Referrals []models.Referral `json:"referrals"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Register User"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Router /register [post]
func (uc *UserController) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := uc.UserService.Register(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login Credentials"
// @Success 200 {object} TokenResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /login [post]
func (uc *UserController) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := uc.UserService.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: token})
}

// CreateReferralCode godoc
// @Summary Create referral code
// @Description Create a new referral code with expiry date
// @Tags referral
// @Accept json
// @Produce json
// @Param referral body CreateReferralRequest true "Referral Code Creation"
// @Success 200 {object} ReferralResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /referral_code [post]
func (uc *UserController) CreateReferralCode(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var req CreateReferralRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := uc.UserService.CreateReferralCode(userID, req.Expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ReferralResponse{
		ReferralCode: user.ReferralCode,
		Expiry:       user.ReferralExpiry,
	})
}

// DeleteReferralCode godoc
// @Summary Delete referral code
// @Description Delete the user's existing referral code
// @Tags referral
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /referral_code [delete]
func (uc *UserController) DeleteReferralCode(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	user, err := uc.UserService.DeleteReferralCode(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message:      "Referral code deleted",
		ReferralCode: user.ReferralCode,
	})
}

// GetReferralCodeByEmail godoc
// @Summary Get referral code by email
// @Description Retrieve a user's referral code using their email
// @Tags referral
// @Produce json
// @Param email query string true "User Email"
// @Success 200 {object} ReferralResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /referral_code [get]
func (uc *UserController) GetReferralCodeByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "email is required"})
		return
	}

	code, err := uc.UserService.GetReferralCodeByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ReferralResponse{ReferralCode: code})
}

// RegisterWithReferral godoc
// @Summary Register with referral code
// @Description Register a new user using a referral code
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterWithReferralRequest true "Register with Referral"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Router /register_with_referral [post]
func (uc *UserController) RegisterWithReferral(c *gin.Context) {
	var req RegisterWithReferralRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := uc.UserService.RegisterWithReferral(req.ReferralCode, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// GetReferrals godoc
// @Summary Get user referrals
// @Description Retrieve a list of users referred by the authenticated user
// @Tags referral
// @Produce json
// @Success 200 {object} ReferralsResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /referrals [get]
func (uc *UserController) GetReferrals(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	referrals, err := uc.UserService.GetReferrals(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ReferralsResponse{Referrals: referrals})
}
