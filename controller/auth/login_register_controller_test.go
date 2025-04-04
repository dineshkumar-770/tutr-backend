package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	controller "main/controller/auth"
	"main/database"
	h "main/helpers"
	"main/models"

	u "main/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateTeacherUserSuccess(t *testing.T) {
	u.GetEnvVars()

	db := database.Initialize()
	if db == nil {
		t.Fatal("Database initialization error:")
	}
	errDbPing := db.Ping()
	if errDbPing != nil {
		t.Fatalf("Failed in database connection error: %v", errDbPing)
	}

	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
	teacherId := h.GenerateUserID("teacher")
	createdAt := time.Now().Unix()
	teacher := models.Teacher{
		Email:           email,
		TeacherID:       teacherId,
		CreatedAt:       createdAt,
		FullName:        "Test Teacher",
		ContactNumber:   987654323,
		Subject:         "Test Subject",
		Qualification:   "Test Qualification",
		ExperienceYears: 12,
		Address:         "Test Address",
		ClassAssigned:   "11TH,12TH",
		TeacherCode:     "TESTCODE",
	}

	body, errM := json.Marshal(teacher)

	if errM != nil {
		t.Fatalf("Unable to parse teacher info. maybe something type mismatch or some wrong info provided: %v", errM)
	}

	req, rqErr := http.NewRequest("POST", "/register_teacher", bytes.NewBuffer(body))
	if rqErr != nil {
		t.Fatalf("Error in requesting URL with body bytes data: %v", rqErr)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller.CreateTeacherUser(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code %v, but got %v. Response: %v", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp u.ResponseStr
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("Expected status 'success', but got '%v'. Response: %v", resp.Status, rr.Body.String())
	}

	// ✅ Step 8: Verify database entry
	defer func() {
		_, err := db.Exec("DELETE FROM register_teachers WHERE email = ?", email)
		if err != nil {
			log.Printf("Failed to delete test entry: %v", err)
		}
	}()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM register_teachers WHERE email = ?", email).Scan(&count)
	fmt.Println("test entries count: ", count)
	if err != nil {
		t.Fatalf("Database query failed: %v", err)
	} else if count != 1 {
		t.Fatalf("Expected 1 database entry, but found %d", count)
	}
}

func TestCreateStudentUserSuccess(t *testing.T) {
	// ✅ Step 1: Load environment variables
	u.GetEnvVars()

	// ✅ Step 2: Initialize database
	db := database.Initialize()
	if db == nil {
		t.Fatal("Database initialization failed")
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("Database connection lost: %v", err)
	}

	// ✅ Step 3: Generate unique email for test
	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())
	student := models.StudentUser{
		FullName:       "Test Student",
		Email:          email,
		Password:       "password123",
		ContactNumber:  1234567890,
		Class:          "10TH",
		ParentsContact: 9876543210,
		FullAddress:    "Test Address",
	}

	// ✅ Step 4: Prepare HTTP request
	body, _ := json.Marshal(student)
	req, _ := http.NewRequest("POST", "/register_student", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// ✅ Step 5: Call controller function
	controller.CreateStudentUser(rr, req)

	// ✅ Step 6: Check HTTP response status
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code %v, but got %v. Response: %v", http.StatusOK, rr.Code, rr.Body.String())
	}

	// ✅ Step 7: Parse response
	var resp u.ResponseStr
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("Expected status 'success', but got '%v'. Response: %v", resp.Status, rr.Body.String())
	}

	// ✅ Step 8: Verify database entry
	defer func() {
		_, err := db.Exec("DELETE FROM register_students WHERE email = ?", email)
		if err != nil {
			log.Printf("Failed to delete test entry: %v", err)
		}
	}()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM register_students WHERE email = ?", email).Scan(&count)
	if err != nil {
		t.Fatalf("Database query failed: %v", err)
	} else if count != 1 {
		t.Fatalf("Expected 1 database entry, but found %d", count)
	}
}
