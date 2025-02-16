package models

type OTPModel struct {
	Otp       string `json:"otp"`
	Expiry    int    `json:"expiry"`
	Email     string `json:"email"`
	UserId    string `json:"user_id"`
	LoginType string `json:"login_type"`
}
