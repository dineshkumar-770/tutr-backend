package middlewares

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("dineshkumar")

func CreateToken(userId string, email string, userType string) (string, error) {
	myToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":   userId,
			"email":     email,
			"user_type": userType,
			"exp":       time.Now().Add(time.Hour * 720).Unix(),
		},
	)

	tokenString, err := myToken.SignedString(secretKey)

	return tokenString, err
}
