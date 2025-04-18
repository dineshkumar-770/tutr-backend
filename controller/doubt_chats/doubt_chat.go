package doubtchats

import (
	"database/sql"
	"fmt"
	"main/constants"
	"main/database"
	"main/helpers"
	doubtchatmodel "main/models/doubt_chat_model"
	u "main/utils"
	"net/http"
	"time"
)

func CreateDoubtChatMessage(w http.ResponseWriter, r *http.Request) {
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

	var doubtDataModel doubtchatmodel.DoubtChatDataModel

	doubtDataModel.DoubtID = helpers.GenerateUserID("doubt")
	doubtDataModel.GroupID = r.FormValue("group_id")
	doubtDataModel.StudentID = r.FormValue("student_id")
	doubtDataModel.TeacherID = r.FormValue("teacher_id")
	doubtDataModel.DoubtText = r.FormValue("text")
	doubtDataModel.FilePath = ""
	doubtDataModel.DoubtStatus = "unsolved"
	doubtDataModel.CreatedAt = int(time.Now().Unix())
	doubtDataModel.UpdatedAt = int(time.Now().Unix())

	if doubtDataModel.GroupID == "" || doubtDataModel.StudentID == "" || doubtDataModel.TeacherID == "" || doubtDataModel.DoubtText == "" {
		u.SendResponseWithMissingValues(w)
		return
	}

	db := database.GetDBInstance()

	_, dbErr := db.Exec("INSERT INTO tutrdevdb.doubts (doubt_id,group_id,student_id,teacher_id,doubt_text,file_url,status,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?)",
		doubtDataModel.DoubtID, doubtDataModel.GroupID, doubtDataModel.StudentID, doubtDataModel.TeacherID, doubtDataModel.DoubtText, doubtDataModel.FilePath, doubtDataModel.DoubtStatus, doubtDataModel.CreatedAt, doubtDataModel.UpdatedAt)

	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			fmt.Println(dbErr)
			resp.MyResponse = dbErr
			resp.Message = "unable to insert data because no entry found in our databse"
			u.SendResponseWithServerError(w, resp)
			return
		}

		fmt.Println(dbErr)
		resp.MyResponse = dbErr
		resp.Message = "unable to insert data. try again later!"
		u.SendResponseWithServerError(w, resp)
		return
	}

	resp.MyResponse = doubtDataModel
	resp.Status = "success"
	resp.Message = "doubt added successfully"
	u.SendResponseWithOK(w, resp)
}

func GetAllDoubtChatsOfGroup(w http.ResponseWriter, r *http.Request) {
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

	groupId := r.FormValue("group_id")
	teacherId := r.FormValue("teacher_id")

	if groupId == "" || teacherId == "" {
		u.SendResponseWithMissingValues(w)
		return
	}

	var doubtChats []doubtchatmodel.DoubtChatData

	db := database.GetDBInstance()

	rows, dbErr := db.Query("SELECT tutrdevdb.doubts.*,tutrdevdb.register_students.full_name,tutrdevdb.register_students.email,tutrdevdb.register_students.class FROM tutrdevdb.doubts RIGHT JOIN tutrdevdb.register_students ON tutrdevdb.register_students.student_id = tutrdevdb.doubts.student_id WHERE tutrdevdb.doubts.group_id = (?) AND tutrdevdb.doubts.teacher_id = (?) ORDER BY tutrdevdb.doubts.created_at DESC;", groupId, teacherId)

	if dbErr != nil {
		resp.Message = "cannot fetch chats right now!"
		fmt.Println(dbErr)
		u.SendResponseWithServerError(w, resp)
		return
	}

	for rows.Next() {
		var doubtChatData doubtchatmodel.DoubtChatData

		err := rows.Scan(&doubtChatData.DoubtID, &doubtChatData.GroupID, &doubtChatData.StudentID, &doubtChatData.TeacherID, &doubtChatData.DoubtText, &doubtChatData.FilePath, &doubtChatData.DoubtStatus, &doubtChatData.CreatedAt, &doubtChatData.UpdatedAt, &doubtChatData.StudentFullName, &doubtChatData.StudentEmail, &doubtChatData.StudentClass)
		if err != nil {
			resp.Message = "cannot read chats of this group"
			fmt.Println(err)
			u.SendResponseWithServerError(w, resp)
			return
		}

		doubtChats = append(doubtChats, doubtChatData)
	}

	if len(doubtChats) == 0 {
		resp.Message = "Type below to share your doubts in group"
		resp.MyResponse = nil
		u.SendResponseWithStatusNotFound(w, resp)
		return
	}

	resp.Status = "success"
	resp.Message = "chats fetched successfully"
	resp.MyResponse = doubtChats
	u.SendResponseWithOK(w, resp)

}
