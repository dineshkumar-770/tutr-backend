package attendance

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tutr-backend/constants"
	"tutr-backend/database"
	"tutr-backend/helpers"
	"tutr-backend/middlewares"
	attendancemodel "tutr-backend/models/attendance_model"
	u "tutr-backend/utils"
)

func MarkStudentAttendance(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		resp.Message = "unable to verify your token"
		fmt.Println(err)
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	teacherId := claims["user_id"].(string)

	if strings.Contains(teacherId, constants.STUDENT) {
		resp.Status = constants.FAILED
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	var attendanceStudent attendancemodel.StudentAttendance
	var attendanceBody attendancemodel.AttendanceBody

	decodeErr := json.NewDecoder(r.Body).Decode(&attendanceBody)
	if decodeErr != nil {
		resp.Message = "Error in parsing sent data. Something unusual if recieved which is not valid. kindly check the field and try again."
		fmt.Println("Decoded error: ", decodeErr)
		resp.Status = constants.FAILED
		u.SendResponseWithStatusBadRequest(w, resp)
		return
	}

	presentStudent := 0
	absentStudent := 0

	for _, val := range attendanceBody.MarkedAttendance {
		switch strings.ToLower(val.Status) {
		case "present":
			presentStudent++
		case "absent":
			absentStudent++
		}
	}

	attendanceID := helpers.GenerateUserID("attendance")

	attendanceStudent.AttendanceID = attendanceID
	attendanceStudent.Remarks = attendanceBody.Remarks
	attendanceStudent.PresentStudents = presentStudent
	attendanceStudent.AbsentStudents = absentStudent
	attendanceStudent.CreatedAt = time.Now().Unix()
	attendanceStudent.IsMarkedSuccessfully = 0
	attendanceStudent.TeacherID = attendanceBody.TeacherID
	attendanceStudent.GroupID = attendanceBody.GroupID
	attendanceStudent.TotalStudents = len(attendanceBody.MarkedAttendance)

	encodedJSON, err := json.Marshal(attendanceBody.MarkedAttendance)
	if err != nil {
		resp.Status = constants.FAILED
		fmt.Println("Error in encoding json: ", err)
		resp.Message = err.Error()
		u.SendResponseWithServerError(w, resp)
		return
	}
	attendanceStudent.MarkedAttendance = string(encodedJSON)

	db := database.GetDBInstance()
	insertQuery := `
	INSERT INTO student_attendance (
		attendance_id, teacher_id, group_id,total_students,
		present_count, absent_count, marked_attendance,
		created_at, remarks, is_marked_successfully
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	_, dbErr := db.Exec(insertQuery,
		attendanceStudent.AttendanceID,
		attendanceStudent.TeacherID,
		attendanceStudent.GroupID,
		attendanceStudent.TotalStudents,
		attendanceStudent.PresentStudents,
		attendanceStudent.AbsentStudents,
		attendanceStudent.MarkedAttendance,
		attendanceStudent.CreatedAt,
		attendanceStudent.Remarks,
		1,
	)

	if dbErr != nil {
		resp.Message = "Failed to save attendance for this group please marks archived and try to save it later!"
		fmt.Println("Insert error: ", dbErr)
		u.SendResponseWithServerError(w, resp)
		return
	}

	attendanceStudent.IsMarkedSuccessfully = 1

	resp.MyResponse = attendanceStudent
	resp.Status = constants.SUCCESS
	resp.Message = "Attendance saved successfully"

	u.SendResponseWithOK(w, resp)

}

func GetAllAttendanceDataOfGroup(w http.ResponseWriter, r *http.Request) {
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
	if err != nil {
		resp.Message = "unable to verify your token"
		fmt.Println(err)
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	teacherId := claims["user_id"].(string)

	if strings.Contains(teacherId, constants.STUDENT) {
		resp.Status = constants.FAILED
		u.SendResponseWithUnauthorizedError(w)
		return
	}

	groupId := r.URL.Query().Get("group_id")
	count := r.URL.Query().Get("count")
	page := r.URL.Query().Get("page")

	if groupId == "" || teacherId == "" {
		resp.Message = "missing field are not allowd"
		resp.Status = constants.FAILED
		u.SendResponseWithStatusNotFound(w, resp)
		return
	}

	if count == "" || page == "" {
		resp.Message = "please specify the page number"
		resp.Status = constants.FAILED
		u.SendResponseWithStatusNotFound(w, resp)
		return
	}

	pageNum, _ := strconv.Atoi(page)
	countNum, _ := strconv.Atoi(count)

	offset := (pageNum - 1) * countNum

	query := `
		SELECT * FROM tutrdevdb.student_attendance WHERE group_id = (?) and teacher_id = (?) ORDER BY created_at DESC LIMIT %d OFFSET %d;
	`

	finalQuery := fmt.Sprintf(query, countNum, offset)

	db := database.GetDBInstance()

	var studentsAttendanceList []attendancemodel.StudentAttendance

	rows, dbErr := db.Query(finalQuery, groupId, teacherId)

	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			resp.Message = "No attendance record found for this group"
			fmt.Println("Insert error:", dbErr)
			u.SendResponseWithStatusNotFound(w, resp)
			return
		}
		resp.Message = "Somthing went wrong occured while getting your attendance record. please try again later"
		fmt.Println("Insert error:", dbErr)
		u.SendResponseWithServerError(w, resp)
		return
	}

	for rows.Next() {
		var attendance attendancemodel.StudentAttendance
		rows.Scan(&attendance.AttendanceID,
			&attendance.TeacherID,
			&attendance.GroupID,
			&attendance.TotalStudents,
			&attendance.PresentStudents,
			&attendance.AbsentStudents,
			&attendance.MarkedAttendance,
			&attendance.CreatedAt,
			&attendance.Remarks,
			&attendance.IsMarkedSuccessfully,
		)

		studentsAttendanceList = append(studentsAttendanceList, attendance)
	}

	if len(studentsAttendanceList) == 0 {
		resp.Message = "No attendance record found for this group"
		u.SendResponseWithStatusNotFound(w, resp)
		return
	}

	resp.MyResponse = studentsAttendanceList
	resp.Message = "Attendance List fetched successfully"
	resp.Status = constants.SUCCESS

	u.SendResponseWithOK(w, resp)

}
