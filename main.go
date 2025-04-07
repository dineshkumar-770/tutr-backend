package main

import (
	"log"
	controller "main/controller/auth"
	doubtchats "main/controller/doubt_chats"
	"main/controller/emails"
	p "main/controller/profile_data"
	s "main/controller/student"
	t "main/controller/teachers"
	tan "main/controller/teachers/add_notes"
	tg "main/controller/teachers/create_group"
	tn "main/controller/teachers/notice_board"
	"main/database"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	database.Initialize()
	r := mux.NewRouter()
	r.HandleFunc("/register_student", controller.CreateStudentUser).Methods("POST")
	r.HandleFunc("/send_otp_by_email", emails.SendOTPByEmail).Methods("POST")
	r.HandleFunc("/verify_otp", emails.VerifyOtp).Methods("POST")
	r.HandleFunc("/register_teacher", controller.CreateTeacherUser).Methods("POST")
	r.HandleFunc("/get_all_teachers", t.GetAvailableTeachers).Methods("GET")
	r.HandleFunc("/create_group", tg.CreateTeacherStudentGroup).Methods("POST")
	r.HandleFunc("/get_teachers_student_groups", tg.GetAllGroupsCreatedByTeacher).Methods("GET")
	r.HandleFunc("/add_student_to_group", tg.AddStudentsToGroup).Methods("POST")
	r.HandleFunc("/check_user_exist", s.CheckStudentExistOrNot).Methods("GET")
	r.HandleFunc("/get_all_group_members_teacher", tg.GetAllStudentsOfGroup).Methods("GET")
	r.HandleFunc("/create_notice_for_group", tn.CreateNoticeForGroup).Methods("POST")
	r.HandleFunc("/update_group_notice", tn.UpdateNoticeBoardForGroup).Methods("POST")
	r.HandleFunc("/get_group_notices", tn.GetGroupNoticeBoard).Methods("POST")
	r.HandleFunc("/insert_doubt_chat", doubtchats.CreateDoubtChatMessage).Methods("POST")
	r.HandleFunc("/get_group_chats", doubtchats.GetAllDoubtChatsOfGroup).Methods("POST")
	r.HandleFunc("/get_user_profile", p.GetUserProfileData).Methods("GET")
	r.HandleFunc("/save_teacher_note", tan.AddTeacherClassNotesToStorage2).Methods("POST")
	r.HandleFunc("/get_all_notes_in_group", tan.GetAllNotesOfTheGroup).Methods("GET")
	log.Println("Server starting on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", r))
}

//Go test run commands
// go test ./controller/... -run "TestSentEmailSuccess|TestCreateTeacherUserSuccess|TestCreateStudentUserSuccess|TestSentEmailTeacherSuccess"
