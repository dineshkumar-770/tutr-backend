package student

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"tutr-backend/constants"
	"tutr-backend/database"
	"tutr-backend/models"
	tg "tutr-backend/models/create_group"
	"tutr-backend/utils"
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

	studentphone := r.URL.Query().Get("phone_number")
	if studentphone == "" {
		utils.SendResponseWithMissingValues(w)
		return
	}

	phoneNum, err := strconv.ParseInt(studentphone, 10, 64)
	if err != nil {
		resp.Message = "invalid phone number format"
		utils.SendResponseWithMissingValues(w)
		return
	}

	db := database.GetDBInstance()
	query := `SELECT * FROM tutrdevdb.register_students WHERE contact_number = ?`
	rows, dbErr := db.Query(query, phoneNum)
	if dbErr != nil {
		resp.Message = "database error"
		resp.MyResponse = dbErr.Error()
		utils.SendResponseWithServerError(w, resp)
		return
	}
	// defer rows.Close()

	var students []models.StudentUserResponse
	for rows.Next() {
		var student models.StudentUserResponse
		err := rows.Scan(
			&student.StudentID,
			&student.FullName,
			&student.Email,
			&student.CreatedAt,
			&student.Password,
			&student.ContactNumber,
			&student.Class,
			&student.TeacherCode,
			&student.ParentsContact,
			&student.FullAddress,
		)
		if err != nil {
			resp.Message = "error scanning student data"
			resp.MyResponse = err.Error()
			utils.SendResponseWithServerError(w, resp)
			return
		}
		students = append(students, student)
	}

	if len(students) == 0 {
		resp.Message = "student not registered with this number"
		utils.SendResponseWithStatusNotFound(w, resp)
		return
	}

	resp.MyResponse = students[0]
	resp.Message = "User found successfully"
	resp.Status = "success"
	utils.SendResponseWithOK(w, resp)
}

func GetAllGroupsWhereStudentAdded(w http.ResponseWriter, r *http.Request, studentId string) {
	resp := utils.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	noTeacherFoundError := "It seems you are not add to any group by any teacher"
	tokenString := r.Header.Get(constants.Authrization)
	if tokenString == "" {
		utils.SendResponseWithUnauthorizedError(w)
		return
	}

	// studentId := r.FormValue("student_id")

	if studentId == "" {
		utils.SendResponseWithMissingValues(w)
		return
	}

	var studentGroupsList []tg.AllTeacherByStudentModel

	db := database.GetDBInstance()

	rows, dbErr := db.Query("SELECT tutrdevdb.teacher_student_group.group_id,tutrdevdb.teacher_student_group.group_name,tutrdevdb.teacher_student_group.group_desc,tutrdevdb.teacher_student_group.created_at,tutrdevdb.teacher_student_group.group_class,tutrdevdb.register_teachers.* FROM tutrdevdb.group_members RIGHT JOIN tutrdevdb.teacher_student_group ON tutrdevdb.group_members.group_id = tutrdevdb.teacher_student_group.group_id RIGHT JOIN  tutrdevdb.register_teachers ON tutrdevdb.teacher_student_group.teacher_id = tutrdevdb.register_teachers.teacher_id WHERE tutrdevdb.group_members.student_id = (?);", studentId)

	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			fmt.Println("parsing student groups ", dbErr)
			resp.Status = "failed"
			resp.Message = noTeacherFoundError
			utils.SendResponseWithStatusNotFound(w, resp)
			return
		}

		fmt.Println("parsing teachers groups ", dbErr)
		resp.Status = "failed"
		resp.Message = "Some error occured while fetching your groups"
		utils.SendResponseWithServerError(w, resp)
		return
	}

	for rows.Next() {
		var studentGroup tg.AllTeacherByStudentModel
		studentGroup.TeacherDetails = &tg.StudentTeacher{}
		err := rows.Scan(
			&studentGroup.GroupId,
			&studentGroup.GroupName,
			&studentGroup.GroupDesctiption,
			&studentGroup.CreatedAt,
			&studentGroup.GroupClass,
			&studentGroup.TeacherDetails.TeacherID,
			&studentGroup.TeacherDetails.FullName,
			&studentGroup.TeacherDetails.Email,
			&studentGroup.TeacherDetails.ContactNumber,
			&studentGroup.TeacherDetails.Subject,
			&studentGroup.TeacherDetails.CreatedAt,
			&studentGroup.TeacherDetails.Qualification,
			&studentGroup.TeacherDetails.ExperienceYears,
			&studentGroup.TeacherDetails.Address,
			&studentGroup.TeacherDetails.ClassAssigned,
			&studentGroup.TeacherDetails.TeacherCode,
		)

		if err != nil {
			resp.Message = "unable to red the groups that you has been added to. Kindly confirm the teacher that you are been added though!"
			fmt.Println(err)
			utils.SendResponseWithServerError(w, resp)
			return
		}

		studentGroupsList = append(studentGroupsList, studentGroup)
	}

	if len(studentGroupsList) == 0 {
		fmt.Println("parsing student groups ", dbErr)
		resp.Status = "failed"
		resp.Message = "It seems you are not add to any group by any teacher"
		utils.SendResponseWithStatusNotFound(w, resp)
		return
	}

	resp.Status = "success"
	resp.Message = "Groups fetched successfully"
	resp.MyResponse = studentGroupsList

	utils.SendResponseWithOK(w, resp)

}
