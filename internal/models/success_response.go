package models

type SuccessResponse struct {
	Message      string `json:"message" example:"Referral code deleted"`
	ReferralCode string `json:"referral_code,omitempty" example:"ABC123XYZ"`
}
