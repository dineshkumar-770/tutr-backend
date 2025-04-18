package noticeboard

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"tutr-backend/database"
	"tutr-backend/helpers"
	tn "tutr-backend/models/notice_board"
	"tutr-backend/utils"
)

func CreateNoticeForGroup(w http.ResponseWriter, r *http.Request) {
	//NoticeBoardModel
	resp := utils.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	tokenData := helpers.VerifyTokenAndGetUserID(w, r)
	if tokenData.UserID == "" || tokenData.UserEmail == "" || tokenData.UserType == "" {
		return
	}

	groupId := r.FormValue("group_id")
	title := r.FormValue("title")
	description := r.FormValue("desc")

	noticeId := helpers.GenerateUserID("notice")
	createdAt := time.Now().Unix()
	updatedAt := time.Now().Unix()

	if groupId == "" || title == "" || description == "" {
		utils.SendResponseWithMissingValues(w)
		return
	}

	db := database.GetDBInstance()
	_, dbErr := db.Exec("INSERT INTO tutrdevdb.group_notice_board (notice_id,title,description,group_id,teacher_id,created_at,updated_at) VALUES (?,?,?,?,?,?,?)",
		noticeId, title, description, groupId, tokenData.UserID, createdAt, updatedAt,
	)

	if dbErr != nil {
		resp.Message = "unable to create notice for this group."
		fmt.Println(dbErr)
		utils.SendResponseWithServerError(w, resp)
		return
	}

	resp.Status = "success"
	resp.Message = "Notice added to your group's notice board"
	resp.MyResponse = nil
	utils.SendResponseWithOK(w, resp)
}

func UpdateNoticeBoardForGroup(w http.ResponseWriter, r *http.Request) {
	resp := utils.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	tokenData := helpers.VerifyTokenAndGetUserID(w, r)
	if tokenData.UserID == "" || tokenData.UserEmail == "" || tokenData.UserType == "" {
		return
	}

	noticeId := r.FormValue("notice_id")
	groupID := r.FormValue("group_id")
	newDesc := r.FormValue("new_desc")
	newTitle := r.FormValue("new_title")

	if noticeId == "" || groupID == "" || newDesc == "" || newTitle == "" {
		utils.SendResponseWithMissingValues(w)
		return
	}

	updateAt := time.Now().Unix()

	db := database.GetDBInstance()
	_, updpateErr := db.Exec("UPDATE tutrdevdb.group_notice_board SET description = (?), title = (?), updated_at = (?) WHERE (notice_id = ? AND group_id = ?)", newDesc, newTitle, updateAt, noticeId, groupID)
	if updpateErr != nil {
		resp.Message = "error in updating current notice"
		fmt.Println(updpateErr)
		utils.SendResponseWithServerError(w, resp)
		return
	}

	var noticeData []tn.NoticeBoardModel

	rows, err := db.Query("SELECT description, title, updated_at, notice_id FROM tutrdevdb.group_notice_board WHERE notice_id = (?) AND group_id = (?)", noticeId, groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			resp.Message = "No notices found for this group"
			fmt.Println(err)
			utils.SendResponseWithServerError(w, resp)
			return
		} else {
			resp.Message = "error in getting current notice"
			fmt.Println(err)
			utils.SendResponseWithServerError(w, resp)
			return
		}
	}

	for rows.Next() {
		var notice tn.NoticeBoardModel
		err := rows.Scan(&notice.Desctiption, &notice.Title, &notice.UpdatedAt, &notice.NoticeID)
		if err != nil {
			fmt.Println("error in getting group notice ", err)
			resp.Status = "failed"
			resp.Message = "Unable to get group notice at this time!"
			utils.SendResponseWithServerError(w, resp)
			return
		}

		noticeData = append(noticeData, notice)
	}

	resp.Status = "success"
	resp.Message = "Notice for this group fetched successfully"
	resp.MyResponse = noticeData
	utils.SendResponseWithOK(w, resp)

}

func GetGroupNoticeBoard(w http.ResponseWriter, r *http.Request) {
	resp := utils.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	tokenData := helpers.VerifyTokenAndGetUserID(w, r)
	if tokenData.UserID == "" || tokenData.UserEmail == "" || tokenData.UserType == "" {
		return
	}

	groupId := r.FormValue("group_id")

	if groupId == "" {
		utils.SendResponseWithMissingValues(w)
		return
	}

	db := database.GetDBInstance()
	var noticeData []tn.NoticeBoardModel

	rows, err := db.Query("SELECT description, title, updated_at, notice_id FROM tutrdevdb.group_notice_board WHERE  group_id = (?)", groupId)
	if err != nil {
		if err == sql.ErrNoRows {
			resp.Message = "No notices found for this group"
			fmt.Println(err)
			utils.SendResponseWithServerError(w, resp)
			return
		} else {
			resp.Message = "error in getting current notice"
			fmt.Println(err)
			utils.SendResponseWithServerError(w, resp)
			return
		}
	}

	fmt.Println("grt notcies rows== ", rows)
	fmt.Println("grt notcies err== ", err)

	for rows.Next() {
		var notice tn.NoticeBoardModel
		err := rows.Scan(&notice.Desctiption, &notice.Title, &notice.UpdatedAt, &notice.NoticeID)
		if err != nil {
			fmt.Println("error in getting group notice ", err)
			resp.Status = "failed"
			resp.Message = "Unable to get group notice at this time!"
			utils.SendResponseWithServerError(w, resp)
			return
		}

		noticeData = append(noticeData, notice)
	}

	if len(noticeData) == 0 {
		resp.Status = "failed"
		resp.Message = "No notices found for this group"
		utils.SendResponseWithOK(w, resp)
		return
	}

	resp.Status = "success"
	resp.Message = "Notice for this group fetched successfully"
	resp.MyResponse = noticeData[0]
	utils.SendResponseWithOK(w, resp)

}
