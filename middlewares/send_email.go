package middlewares

import (
	"crypto/rand"
	"fmt"

	u "main/utils"
	"math/big"
	"net/smtp"
)

func getHTMLOtpText(otp string) string {
	var contentHtml = "Your OTP is %s"
	return fmt.Sprintf(contentHtml, otp)
}

func SendEmail(emailId string) (string, error) {
	envs, err := u.GetEnvVars()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	sentOTP := generateOTP(6)

	htmlText := getHTMLOtpText(sentOTP)
	listToSentEmail := []string{emailId}
	host := "smtp.gmail.com"
	port := "587"
	body := []byte(htmlText)
	auth := smtp.PlainAuth("", envs.AppEmailID, envs.EmailAppPassword, host)
	err1 := smtp.SendMail(host+":"+port, auth, envs.AppEmailID, listToSentEmail, body)
	if err1 != nil {
		fmt.Println(err)
		return "", err1
	}

	fmt.Println("Successfully sent mail to all user in toList")
	fmt.Println(sentOTP)
	return sentOTP, nil
}

func generateOTP(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		otp[i] = charset[num.Int64()]
	}
	return string(otp)
}
