package attendancemodel

type StudentAttendance struct {
	AttendanceID         string `json:"attendance_id,omitempty"`
	TeacherID            string `json:"teacher_id,omitempty"`
	GroupID              string `json:"group_id,omitempty"`
	PresentStudents      int    `json:"present_student,omitempty"`
	TotalStudents        int    `json:"total_students,omitempty"`
	AbsentStudents       int    `json:"absent_students,omitempty"`
	MarkedAttendance     string `json:"marked_attendance,omitempty"`
	CreatedAt            int64  `json:"created_at,omitempty"`
	Remarks              string `json:"remarks,omitempty"`
	IsMarkedSuccessfully int    `json:"is_marked_success,omitempty"`
}

type AttendanceBody struct {
	Remarks          string            `json:"remarks,omitempty"`
	TeacherID        string            `json:"teacher_id,omitempty"`
	GroupID          string            `json:"group_id,omitempty"`
	MarkedAttendance []AttendanceValue `json:"marked_attendance,omitempty"`
}

type AttendanceValue struct {
	StudentID    string `json:"student_id"`
	Status       string `json:"status"`
	StudentName  string `json:"name"`
	StudentEmail string `json:"st_email"`
	StudentPhone int64  `json:"student_phone"`
}
