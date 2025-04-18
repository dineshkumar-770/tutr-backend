package middlewares

import (
	"crypto/rand"
	"fmt"

	"math/big"
	"net/smtp"
	u "tutr-backend/utils"
)

func getHTMLOtpText(otp string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>TuTr OTP Verification</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0;">
			<div style="max-width: 600px; margin: auto; padding: 20px; background-color: #ffffff; border-radius: 10px; box-shadow: 0 2px 5px rgba(0,0,0,0.1);">
				<div style="text-align: center; margin-bottom: 30px;">
					<h2 style="color: #333;">Welcome to <span style="color: #007bff;">TuTr</span></h2>
					<p style="font-size: 16px; color: #666;">Your Tuition Management Partner</p>
				</div>

				<div style="text-align: center; padding: 20px 0;">
					<p style="font-size: 18px; color: #333;">Use the OTP below to verify your email address:</p>
					<h1 style="font-size: 40px; color: #007bff; margin: 20px 0;">%s</h1>
					<p style="font-size: 16px; color: #999;">This OTP is valid for the next 15 minutes.</p>
				</div>

				<hr style="border: none; height: 1px; background-color: #eaeaea; margin: 30px 0;">

				<div style="text-align: center;">
					<p style="font-size: 14px; color: #aaa;">If you didn’t request this, please ignore this email.</p>
					<p style="font-size: 14px; color: #aaa;">&copy; 2025 TuTr. All rights reserved.</p>
				</div>
			</div>
		</body>
		</html>
	`, otp)
}
func SendEmail(emailId string) (string, error) {
	envs, err := u.GetEnvVars()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	sentOTP := generateOTP(6)
	htmlText := getHTMLOtpText(sentOTP)

	subject := "Subject: Your OTP Code for TuTr\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime + htmlText)

	from := envs.AppEmailID
	password := envs.EmailAppPassword
	to := []string{emailId}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("Error while sending:", err)
		return "", err
	}

	fmt.Println("✅ Email sent successfully to:", emailId)
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
