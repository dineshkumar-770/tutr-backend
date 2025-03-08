package controller

import (
	"encoding/json"
	"fmt"
	"main/database"
	h "main/helpers"
	"main/models"
	u "main/utils"
	"net/http"
	"time"
)

func CreateStudentUser(w http.ResponseWriter, r *http.Request) {
	mr := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	var student models.StudentUser

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		fmt.Println(err)
		mr.Message = "missing field not allowd"
		u.SendResponseWithStatusBadRequest(w, mr)
		return
	}
	student.StudentID = h.GenerateUserID("student")
	student.CreatedAt = int(time.Now().Unix())
	student.TeacherCode = ""
	db := database.GetDBInstance()
	_, errdb := db.Exec("INSERT INTO register_students (student_id,full_name,email,created_at,password,contact_number,class,teacher_code,parents_contact,full_address) VALUES (?,?,?,?,?,?,?,?,?,?)",
		student.StudentID, student.FullName, student.Email, student.CreatedAt, student.Password, student.ContactNumber, student.Class, student.TeacherCode, student.ParentsContact, student.FullAddress,
	)

	if errdb != nil {
		fmt.Println(errdb)
		mr.MyResponse = errdb
		mr.Message = "cannot save user info or cannot create user at this time"
		u.SendResponseWithServerError(w, mr)
		return
	}

	mr.MyResponse = student
	mr.Status = "success"
	mr.Message = fmt.Sprintf("Student registerd successfully with student id %s", student.StudentID)
	u.SendResponseWithOK(w, mr)

}

func CreateTeacherUser(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}
	var teacher models.Teacher
	err := json.NewDecoder(r.Body).Decode(&teacher)
	if err != nil {
		fmt.Println(err)
		resp.MyResponse = err
		resp.Message = "missing field not allowd"
		u.SendResponseWithStatusBadRequest(w, resp)
		return
	}
	teacher.TeacherID = h.GenerateUserID("teacher")
	teacher.CreatedAt = time.Now().Unix()

	db := database.GetDBInstance()
	_, errdb := db.Exec("INSERT INTO register_teachers (teacher_id,full_name,email,contact_number,subject,created_at,qualification,experience_years,address,class_assigned,teacher_code) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		teacher.TeacherID, teacher.FullName, teacher.Email, teacher.ContactNumber, teacher.Subject, teacher.CreatedAt, teacher.Qualification, teacher.ExperienceYears, teacher.Address, teacher.ClassAssigned, teacher.TeacherCode,
	) 
	if errdb != nil {
		fmt.Println(errdb)
		resp.Message = "cannot save user info or cannot create user at this time"
		u.SendResponseWithServerError(w, resp)
		return
	}

	resp.Status = "success"
	resp.Message = fmt.Sprintf("Teacher registered successfully with teacher id %s", teacher.TeacherID)
	u.SendResponseWithOK(w, resp)
}
