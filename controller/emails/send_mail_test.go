package emails_test

import (
	"fmt"
	"main/controller/emails"
	"main/database"
	u "main/utils"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSentEmailStudentSuccess(t *testing.T) {
	u.GetEnvVars()

	db := database.Initialize()
	if db == nil {
		t.Fatal("Database initialization error:")
	}
	errDbPing := db.Ping()
	if errDbPing != nil {
		t.Fatalf("Failed in database connection error: %v", errDbPing)
	}

	//Use the existing emails here or else later will add to register email too and test that email
	email := "ashusharma@gmail.com"
	loginType := "student"

	formValue := url.Values{}
	formValue.Set("email", email)
	formValue.Set("login_type", loginType)

	req, rqErr := http.NewRequest("POST", "send_otp_by_email", strings.NewReader(formValue.Encode()))

	if rqErr != nil {
		t.Fatalf("Failed to create request: %v", rqErr)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	emails.SendOTPByEmail(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	// ✅ Step 6: Verify OTP is Stored in DB
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user_otps WHERE email = ?", email).Scan(&count)
	if err != nil {
		t.Fatalf("Database query failed: %v", err)
	}
	if count == 0 {
		t.Fatalf("Expected OTP entry in DB but found none")
	}

	fmt.Println("Counter: ", count)

	// ✅ Step 7: Cleanup - Delete OTP Entry After Test
	defer func() {
		_, _ = db.Exec("DELETE FROM user_otps WHERE email = ?", email)
		fmt.Println("Final Counter: ", count)
	}()

}
func TestSentEmailTeacherSuccess(t *testing.T) {
	u.GetEnvVars()

	db := database.Initialize()
	if db == nil {
		t.Fatal("Database initialization error:")
	}
	errDbPing := db.Ping()
	if errDbPing != nil {
		t.Fatalf("Failed in database connection error: %v", errDbPing)
	}

	//Use the existing emails here or else later will add to register email too and test that email
	email := "vijaygargmaths@gmail.com"
	loginType := "teacher"

	formValue := url.Values{}
	formValue.Set("email", email)
	formValue.Set("login_type", loginType)

	req, rqErr := http.NewRequest("POST", "send_otp_by_email", strings.NewReader(formValue.Encode()))

	if rqErr != nil {
		t.Fatalf("Failed to create request: %v", rqErr)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	emails.SendOTPByEmail(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Response: %s", rr.Code, rr.Body.String())
	}

	// ✅ Step 6: Verify OTP is Stored in DB
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user_otps WHERE email = ?", email).Scan(&count)
	if err != nil {
		t.Fatalf("Database query failed: %v", err)
	}
	if count == 0 {
		t.Fatalf("Expected OTP entry in DB but found none")
	}

	fmt.Println("Counter: ", count)

	// ✅ Step 7: Cleanup - Delete OTP Entry After Test
	defer func() {
		_, _ = db.Exec("DELETE FROM user_otps WHERE email = ?", email)
		fmt.Println("Final Counter: ", count)
	}()

}
