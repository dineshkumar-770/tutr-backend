package emails

import (
	"database/sql"
	"fmt"
	d "main/database"
	"main/middlewares"
	"main/models"
	u "main/utils"
	"net/http"
	"strings"
	"time"
)

func SendOTPByEmail(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "",
		MyResponse: nil,
	}
	email := r.FormValue("email")
	loginType := r.FormValue("login_type")

	if email == "" || loginType == "" {
		resp.Message = "missing field email is required!"
		u.SendResponseWithStatusNotFound(w, resp)
		return
	}

	var studentID string
	var teacherID string
	otpExpiry := time.Now().Add(time.Minute * 10).Unix()
	db := d.GetDBInstance()
	if loginType == "student" {
		stDBErr := db.QueryRow("SELECT student_id FROM register_students WHERE email = (?)", email).Scan(&studentID)
		fmt.Println(studentID)
		if stDBErr != nil {
			if stDBErr == sql.ErrNoRows {
				resp.Message = "User not found"
				u.SendResponseWithStatusNotFound(w, resp)
				return
			} else {
				u.SendResponseWithServerError(w, resp)
				return
			}
		}

		otp, err := middlewares.SendEmail(email)
		if err != nil {
			resp.Message = "unable to send otp!"
			u.SendResponseWithStatusBadRequest(w, resp)
			return
		}

		db.Exec("INSERT INTO user_otps (otp,expiry,email,user_id,login_type) VALUES (?,?,?,?,?)", otp, otpExpiry, email, studentID, loginType)
		resp.Status = "success"
		resp.Message = "OTP sent successfully"
		u.SendResponseWithOK(w, resp)
		return
	} else if loginType == "teacher" {
		tDBErr := db.QueryRow("SELECT teacher_id FROM register_teachers WHERE email = (?)", email).Scan(&teacherID)
		fmt.Println("teacher id ", teacherID)
		if tDBErr != nil {
			if tDBErr == sql.ErrNoRows {
				resp.Message = "User not found"
				u.SendResponseWithStatusNotFound(w, resp)
				return
			} else {
				u.SendResponseWithServerError(w, resp)
				return
			}
		}

		otp, err := middlewares.SendEmail(email)
		if err != nil {
			resp.Message = "unable to send otp!"
			u.SendResponseWithStatusBadRequest(w, resp)
			return
		}

		db.Exec("INSERT INTO user_otps (otp,expiry,email,user_id,login_type) VALUES (?,?,?,?,?)", otp, otpExpiry, email, teacherID, loginType)
		resp.Status = "success"
		resp.Message = "OTP sent successfully"
		u.SendResponseWithOK(w, resp)
		return
	}
}

func VerifyOtp(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "",
		MyResponse: nil,
	}
	submittedOTP := r.FormValue("otp")
	emailId := r.FormValue("email")
	loginType := r.FormValue("login_type")

	if submittedOTP == "" || emailId == "" || loginType == "" {
		resp.Message = "missing field email is required!"
		u.SendResponseWithStatusNotFound(w, resp)
		return
	}

	var otpUser models.OTPModel
	db := d.GetDBInstance()
	dbErr := db.QueryRow("SELECT otp, expiry, email,user_id,login_type FROM user_otps WHERE email = (?) AND otp = (?) AND login_type = (?)", emailId, submittedOTP, loginType).Scan(
		&otpUser.Otp, &otpUser.Expiry, &otpUser.Email, &otpUser.UserId, &otpUser.LoginType)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			resp.Message = "Invalid OTP!"
			resp.MyResponse = dbErr
			u.SendResponseWithStatusNotFound(w, resp)
			return
		} else {
			resp.MyResponse = dbErr
			u.SendResponseWithServerError(w, resp)
			return
		}
	}
	if strings.Contains(strings.ToLower(otpUser.Otp), strings.ToLower(submittedOTP)) {
		otpGeneratedAt := int64(otpUser.Expiry)
		expiryDuration := int64(15)
		time.Sleep(2 * time.Second)
		if !isOtpExpired(otpGeneratedAt, expiryDuration) {
			resp.Message = "OTP verified successfully"
			resp.Status = "success"
			token, err := middlewares.CreateToken(otpUser.UserId, otpUser.Email, otpUser.LoginType)
			if err != nil {
				fmt.Println("token error: ", err)
			}

			resp.MyResponse = token
			u.SendResponseWithOK(w, resp)
			return
		} else {
			resp.Message = "OTP expired"
			u.SendResponseWithStatusBadRequest(w, resp)
			return
		}

	} else {
		resp.Message = "Invalid OTP !"
		u.SendResponseWithStatusBadRequest(w, resp)
		return
	}
}

func isOtpExpired(otpTimestamp int64, expiryDurationInSeconds int64) bool {
	currentTime := time.Now().Unix()
	expiryTime := otpTimestamp + expiryDurationInSeconds
	return currentTime > expiryTime
}
