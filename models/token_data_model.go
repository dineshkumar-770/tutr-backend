package models

type TokenDataModel struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"email"`
	UserType  string `json:"user_type"`
}
