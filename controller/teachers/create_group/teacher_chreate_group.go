package creategroup

import (
	"encoding/json"
	"fmt"
	"main/constants"
	"main/database"
	"main/helpers"
	"main/middlewares"
	"main/models"
	tg "main/models/create_group"
	u "main/utils"
	"net/http"
	"time"
)

func CreateTeacherStudentGroup(w http.ResponseWriter, r *http.Request) {
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
		u.SendResponseWithServerError(w, resp)
		return
	}
	groupClass := r.FormValue("group_class")
	groupID := helpers.GenerateUserID("group")
	groupName := r.FormValue("group_name")
	groupDescription := r.FormValue("group_desc")
	teacherID := claims["user_id"]
	createdAt := time.Now().Unix()

	if groupClass == "" || groupID == "" || groupName == "" {
		resp.Message = "missing field required "
		u.SendResponseWithStatusBadRequest(w, resp)
		return
	}

	db := database.GetDBInstance()
	_, dbErr := db.Exec("INSERT INTO teacher_student_group (group_id,group_name,teacher_id,created_at,group_class,group_desc) VALUES (?,?,?,?,?,?);", groupID, groupName, teacherID, createdAt, groupClass, groupDescription)
	if dbErr != nil {
		resp.Message = "unable to create group at this time"
		fmt.Println(dbErr)
		u.SendResponseWithServerError(w, resp)
		return
	}

	resp.Message = "group created successfully"
	resp.Status = "success"
	u.SendResponseWithOK(w, resp)

}

func GetAllGroupsCreatedByTeacher(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}
	var teacherGroups []tg.CreateGroupModel

	tokenData := helpers.VerifyTokenAndGetUserID(w, r)
	if tokenData.UserID == "" || tokenData.UserEmail == "" || tokenData.UserType == "" {
		return
	}

	db := database.GetDBInstance()
	rows, dbErr := db.Query("SELECT * FROM tution_management.teacher_student_group WHERE teacher_id = (?) ORDER BY created_at DESC", tokenData.UserID)

	if dbErr != nil {
		fmt.Println("parsing teachers groups ", dbErr)
		resp.Status = "failed"
		resp.Message = "could not found any group of yours!"
		u.SendResponseWithServerError(w, resp)
		return
	}

	for rows.Next() {
		var group tg.CreateGroupModel

		var members []tg.GroupMemberStudentsData
		err := rows.Scan(&group.GroupID, &group.GroupName, &group.TeacherID, &group.CreatedAt, &group.GroupClass, &group.GroupDescription)
		if err != nil {
			fmt.Println("parsing teacher groups ", err)
			resp.Status = "failed"
			resp.Message = "could not found any group of yours!"
			u.SendResponseWithServerError(w, resp)
			return
		}

		rows1, memberErr := db.Query("select * from tution_management.group_members where group_id = (?) and owner_id = (?)", group.GroupID, tokenData.UserID)
		for rows1.Next() {
			var member tg.GroupMemberStudentsData
			err := rows1.Scan(&member.GroupMemberID, &member.GroupID, &member.StudentID, &member.OwnerID, &member.GroupOwner, &member.StudentJoinedAt, &member.StudentFullName, &member.StudentEmail, &member.StudentAccountCreatedAt, &member.StudentContactNumber, &member.StudentClass, &member.StudentParentsContact, &member.StudentFullAddress)
			if err != nil {
				resp.Message = "cannot read members of this group"
				fmt.Println(err)
				u.SendResponseWithServerError(w, resp)
				return
			}

			members = append(members, member)
		}

		group.AllMembers = &members

		if memberErr != nil {
			fmt.Println("fetching members from db error: ", memberErr)
		}

		teacherGroups = append(teacherGroups, group)
	}

	resp.Status = "success"
	resp.Message = "groups fetched successfully"
	resp.MyResponse = teacherGroups
	u.SendResponseWithOK(w, resp)

}

func AddStudentsToGroup(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}
	var studentData models.StudentModelAddGroup
	json.NewDecoder(r.Body).Decode(&studentData)
	memberId := helpers.GenerateUserID("member")
	joinedAt := time.Now().Unix()

	if studentData.GroupID == "" || studentData.FullName == "" {
		resp.Message = "missing field not allowd"
		u.SendResponseWithStatusBadRequest(w, resp)
		return
	}
	tokenData := helpers.VerifyTokenAndGetUserID(w, r)
	if tokenData.UserID == "" || tokenData.UserEmail == "" || tokenData.UserType == "" {
		return
	}

	db := database.GetDBInstance()
	_, dbErr := db.Exec("INSERT INTO tution_management.group_members (group_member_id,group_id,student_id,owner_id,group_owner,joined_at,full_name,email,created_at,contact_number,class,parents_contact,full_address) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)",
		memberId, studentData.GroupID, studentData.StudentID, tokenData.UserID, studentData.OwnerName, joinedAt, studentData.FullName, studentData.Email, studentData.CreatedAt, studentData.ContactNumber, studentData.Class, studentData.ParentsContact, studentData.FullAddress)
	if dbErr != nil {
		resp.Message = "cannot add to the group"
		fmt.Println(dbErr)
		u.SendResponseWithServerError(w, resp)
		return
	}

	resp.Message = "Student added to your group successfully"
	resp.Status = "success"
	resp.MyResponse = nil
	u.SendResponseWithOK(w, resp)
}

func GetAllStudentsOfGroup(w http.ResponseWriter, r *http.Request) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}

	tokenData := helpers.VerifyTokenAndGetUserID(w, r)
	if tokenData.UserID == "" || tokenData.UserEmail == "" || tokenData.UserType == "" {
		return
	}

	var members []tg.GroupMemberStudentsData

	ownerID := tokenData.UserID
	groupID := r.FormValue("group_id")

	if groupID == "" {
		u.SendResponseWithMissingValues(w)
		return
	}

	db := database.GetDBInstance()

	rows, dbErr := db.Query("select * from tution_management.group_members where group_id = (?) and owner_id = (?)", groupID, ownerID)
	if dbErr != nil {
		resp.Message = "cannot fetch member of this group at now"
		fmt.Println(dbErr)
		u.SendResponseWithServerError(w, resp)
		return
	}

	for rows.Next() {
		var member tg.GroupMemberStudentsData
		err := rows.Scan(&member.GroupMemberID, &member.GroupID, &member.StudentID, &member.OwnerID, &member.GroupOwner, &member.StudentJoinedAt, &member.StudentFullName, &member.StudentEmail, &member.StudentAccountCreatedAt, &member.StudentContactNumber, &member.StudentClass, &member.StudentParentsContact, &member.StudentFullAddress)
		if err != nil {
			resp.Message = "cannot read members of this group"
			fmt.Println(err)
			u.SendResponseWithServerError(w, resp)
			return
		}

		members = append(members, member)
	}

	resp.Status = "success"
	resp.Message = "member fetched successfully"
	resp.MyResponse = members
	u.SendResponseWithOK(w, resp)

}
