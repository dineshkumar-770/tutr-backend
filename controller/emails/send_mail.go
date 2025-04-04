package emails

import (
	"database/sql"
	"fmt"
	"main/constants"
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
		Status:     constants.FAILED,
		Message:    "",
		MyResponse: nil,
	}
	email := r.FormValue("email")
	loginType := r.FormValue("login_type")

	if email == "" || loginType == "" {
		resp.Message = constants.MISSING_FIELD_EMAIL_MESSAGE
		u.SendResponseWithStatusNotFound(w, resp)
		return
	}

	var studentID string
	var teacherID string
	otpExpiry := time.Now().Add(time.Minute * 10).Unix()
	db := d.GetDBInstance()
	if loginType == constants.STUDENT {
		stDBErr := db.QueryRow("SELECT student_id FROM register_students WHERE email = (?)", email).Scan(&studentID)
		fmt.Println(studentID)
		if stDBErr != nil {
			if stDBErr == sql.ErrNoRows {
				resp.Message = constants.USER_NOT_REGISTERED_MESSAGE
				u.SendResponseWithStatusNotFound(w, resp)
				return
			} else {
				u.SendResponseWithServerError(w, resp)
				return
			}
		}

		otp, err := middlewares.SendEmail(email)
		if err != nil {
			resp.Message = constants.OTP_SENT_ERROR
			u.SendResponseWithStatusBadRequest(w, resp)
			return
		}

		db.Exec("INSERT INTO user_otps (otp,expiry,email,user_id,login_type) VALUES (?,?,?,?,?)", otp, otpExpiry, email, studentID, loginType)
		resp.Status = constants.SUCCESS
		resp.Message = constants.OTP_SENT_SUCCESS
		resp.MyResponse = otp
		u.SendResponseWithOK(w, resp)
		return
	} else if loginType == constants.TEACHER {
		tDBErr := db.QueryRow("SELECT teacher_id FROM register_teachers WHERE email = (?)", email).Scan(&teacherID)
		fmt.Println("teacher id ", teacherID)
		if tDBErr != nil {
			if tDBErr == sql.ErrNoRows {
				resp.Message = constants.USER_NOT_REGISTERED_MESSAGE
				u.SendResponseWithStatusNotFound(w, resp)
				return
			} else {
				u.SendResponseWithServerError(w, resp)
				return
			}
		}

		otp, err := middlewares.SendEmail(email)
		if err != nil {
			resp.Message = constants.OTP_SENT_ERROR
			u.SendResponseWithStatusBadRequest(w, resp)
			return
		}

		_, exeErr := db.Exec("INSERT INTO user_otps (otp,expiry,email,user_id,login_type) VALUES (?,?,?,?,?)", otp, otpExpiry, email, teacherID, loginType)
		if exeErr != nil {
			resp.Message = constants.DB_INSERT_OTP_FAILED
			fmt.Println(exeErr)
			u.SendResponseWithServerError(w, resp)
			return
		}
		resp.Status = constants.SUCCESS
		resp.Message = constants.OTP_SENT_SUCCESS
		fmt.Println("otp request my user", otp)
		resp.MyResponse = otp
		u.SendResponseWithOK(w, resp)
		return
	}
}

func VerifyOtp(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     constants.FAILED,
		Message:    "",
		MyResponse: nil,
	}
	submittedOTP := r.FormValue("otp")
	emailId := r.FormValue("email")
	loginType := r.FormValue("login_type")

	if submittedOTP == "" || emailId == "" || loginType == "" {
		resp.Message = constants.MISSING_FIELD_EMAIL_MESSAGE
		u.SendResponseWithStatusNotFound(w, resp)
		return
	}

	var otpUser models.OTPModel
	db := d.GetDBInstance()
	dbErr := db.QueryRow("SELECT otp, expiry, email,user_id,login_type FROM user_otps WHERE email = (?) AND otp = (?) AND login_type = (?)", emailId, submittedOTP, loginType).Scan(
		&otpUser.Otp, &otpUser.Expiry, &otpUser.Email, &otpUser.UserId, &otpUser.LoginType)
	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			resp.Message = constants.INVALID_OTP_MESSAGE
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
			resp.Message = constants.OTP_VERIFIED_MESSAGE
			resp.Status = constants.SUCCESS
			token, err := middlewares.CreateToken(otpUser.UserId, otpUser.Email, otpUser.LoginType)
			if err != nil {
				fmt.Println("token error: ", err)
			}

			resp.MyResponse = token
			u.SendResponseWithOK(w, resp)
			return
		} else {
			resp.Message = constants.OTP_EXPIRED_MESSAGE
			u.SendResponseWithStatusBadRequest(w, resp)
			return
		}

	} else {
		resp.Message = constants.INVALID_OTP_MESSAGE
		u.SendResponseWithStatusBadRequest(w, resp)
		return
	}
}

func isOtpExpired(otpTimestamp int64, expiryDurationInSeconds int64) bool {
	currentTime := time.Now().Unix()
	expiryTime := otpTimestamp + expiryDurationInSeconds
	return currentTime > expiryTime
}
