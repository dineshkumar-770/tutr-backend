package teachers

import (
	"database/sql"
	"fmt"
	"net/http"
	"tutr-backend/database"
	"tutr-backend/models"
	"tutr-backend/utils"
)

func GetAvailableTeachers(w http.ResponseWriter, r *http.Request) {
	resp := utils.ResponseStr{
		Status:     "failed",
		Message:    "",
		MyResponse: nil,
	}

	var allTeachers []models.Teacher

	db := database.GetDBInstance()
	rows, err := db.Query("SELECT * FROM register_teachers ORDER BY full_name ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			resp.Status = "failed"
			resp.Message = "No teachers available!"
			utils.SendResponseWithStatusNotFound(w, resp)
			return
		} else {
			resp.Status = "failed"
			resp.Message = "Unable to get teachers at this time!"
			fmt.Println("all teachers data ", err)
			utils.SendResponseWithServerError(w, resp)
			return
		}
	}

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.TeacherID, &teacher.FullName, &teacher.Email, &teacher.ContactNumber, &teacher.Subject, &teacher.CreatedAt, &teacher.Qualification, &teacher.ExperienceYears, &teacher.Address, &teacher.ClassAssigned, &teacher.TeacherCode)
		if err != nil {
			fmt.Println("parsing teachers list ", err)
			resp.Status = "failed"
			resp.Message = "Unable to get teachers at this time!"
			utils.SendResponseWithServerError(w, resp)
			return
		}

		allTeachers = append(allTeachers, teacher)
	}

	resp.Status = "success"
	resp.Message = fmt.Sprintf("%d Teachers are available", len(allTeachers))
	resp.MyResponse = allTeachers
	utils.SendResponseWithOK(w, resp)
}
