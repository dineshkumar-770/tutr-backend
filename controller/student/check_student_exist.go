package student

import (
	"fmt"
	"main/constants"
	"main/database"
	"main/models"
	"main/utils"
	"net/http"
	"strconv"
)

func CheckStudentExistOrNot(w http.ResponseWriter, r *http.Request) {
	resp := utils.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}
	tokenString := r.Header.Get(constants.Authrization)
	if tokenString == "" {
		utils.SendResponseWithUnauthorizedError(w)
		return
	}

	student_phone := r.FormValue("phone_number")

	if student_phone == "" {
		utils.SendResponseWithMissingValues(w)
		return
	}

	var students []models.StudentUserResponse

	db := database.GetDBInstance()
	phoneNum, _ := strconv.ParseInt(student_phone, 10, 64)
	rows, dbErr := db.Query("SELECT * FROM tution_management.register_students WHERE contact_number = ?", phoneNum)
	if dbErr != nil {
		resp.Message = "student not registered with this number"
		utils.SendResponseWithStatusNotFound(w, resp)
		return
	}

	for rows.Next() {
		var student models.StudentUserResponse
		rowErr := rows.Scan(&student.StudentID, &student.FullName, &student.Email, &student.CreatedAt, &student.Password, &student.ContactNumber, &student.Class, &student.TeacherCode, &student.ParentsContact, &student.FullAddress)
		if rowErr != nil{
			resp.MyResponse = rowErr
			utils.SendResponseWithServerError(w,resp)
			return
		}
		fmt.Println(student)
		students = append(students, student)
	}

	resp.MyResponse = students[0]
	resp.Message = "User added successfully"
	resp.Status = "success"
	utils.SendResponseWithOK(w, resp)

}
