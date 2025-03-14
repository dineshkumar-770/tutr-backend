package main

import (
	"log"
	controller "main/controller/auth"
	doubtchats "main/controller/doubt_chats"
	"main/controller/emails"
	s "main/controller/student"
	t "main/controller/teachers"
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
	r.HandleFunc("/get_teachers_groups", tg.GetAllGroupsCreatedByTeacher).Methods("GET")
	r.HandleFunc("/add_student_to_group", tg.AddStudentsToGroup).Methods("POST")
	r.HandleFunc("/check_user_exist", s.CheckStudentExistOrNot).Methods("POST")
	r.HandleFunc("/get_student_groups", s.GetAllGroupsWhereStudentAdded).Methods("POST")
	r.HandleFunc("/get_all_group_members_teacher", tg.GetAllStudentsOfGroup).Methods("POST")
	r.HandleFunc("/create_notice_for_group", tn.CreateNoticeForGroup).Methods("POST")
	r.HandleFunc("/update_group_notice", tn.UpdateNoticeBoardForGroup).Methods("POST")
	r.HandleFunc("/get_group_notices", tn.GetGroupNoticeBoard).Methods("POST")
	r.HandleFunc("/insert_doubt_chat", doubtchats.CreateDoubtChatMessage).Methods("POST")
	r.HandleFunc("/get_group_chats", doubtchats.GetAllDoubtChatsOfGroup).Methods("POST")
	log.Println("Server starting on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", r))
}
