package profiledata

import (
	"database/sql"
	"fmt"
	"main/constants"
	"main/database"
	"main/middlewares"
	"main/models"
	u "main/utils"
	"net/http"
	"strings"
)

func GetUserProfileData(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}
	tokenString := r.Header.Get(constants.Authrization)
	if tokenString == "" {
		u.SendResponseWithUnauthorizedError(w)
		return
	}
	tokenString = tokenString[len("Bearer "):]
	claims, err := middlewares.VerifyToken(tokenString)
	if err != nil || claims["user_id"] == ""{
		resp.Message = "unable to verify your token"
		fmt.Println(err)
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	userId := claims["user_id"].(string)

	if userId == "" {
		resp.Message = "unable to verify your token"
		fmt.Println(err)
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	db := database.GetDBInstance()

	if strings.Contains(userId, "student") {
		var studentData models.StudentUserResponse
		dbErr := db.QueryRow("SELECT * FROM tutrdevdb.register_students WHERE student_id = (?)", userId).Scan(
			&studentData.StudentID,
			&studentData.FullName,
			&studentData.Email,
			&studentData.CreatedAt,
			&studentData.Password,
			&studentData.ContactNumber,
			&studentData.Class,
			&studentData.TeacherCode,
			&studentData.ParentsContact,
			&studentData.FullAddress,
		)
		if dbErr != nil {
			if dbErr == sql.ErrNoRows {
				resp.Message = "No user found associated with this user id."
				fmt.Println(dbErr)
				u.SendResponseWithStatusNotFound(w, resp)
				return
			}

			resp.Message = "unable to fetch user profile at this time."
			fmt.Println(dbErr)
			u.SendResponseWithServerError(w, resp)
			return
		}

		resp.Status = "success"
		resp.Message = "Student profile fetched successfully"
		resp.MyResponse = studentData
		u.SendResponseWithOK(w, resp)
		return
	} else {
		var teacherData models.Teacher
		dbErr := db.QueryRow("SELECT * FROM tutrdevdb.register_teachers WHERE register_teachers.teacher_id = (?)", userId).Scan(
			&teacherData.TeacherID,
			&teacherData.FullName,
			&teacherData.Email,
			&teacherData.ContactNumber,
			&teacherData.Subject,
			&teacherData.CreatedAt,
			&teacherData.Qualification,
			&teacherData.ExperienceYears,
			&teacherData.Address,
			&teacherData.ClassAssigned,
			&teacherData.TeacherCode,
		)

		if dbErr != nil {
			if dbErr == sql.ErrNoRows {
				resp.Message = "No user found associated with this user id."
				fmt.Println(dbErr)
				u.SendResponseWithStatusNotFound(w, resp)
				return
			}

			resp.Message = "unable to fetch user profile at this time."
			fmt.Println(dbErr)
			u.SendResponseWithServerError(w, resp)
			return
		}

		resp.Status = "success"
		resp.Message = "Teacher profile fetched successfully"
		resp.MyResponse = teacherData
		u.SendResponseWithOK(w, resp)
		return
	}
}
