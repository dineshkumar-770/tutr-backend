package models

type StudentUser struct {
	StudentID      string `json:"student_id"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	CreatedAt      int    `json:"created_at"`
	Password       string `json:"password"`
	ContactNumber  int64  `json:"contact_number"`
	Class          string `json:"class"`
	TeacherCode    string `json:"teacher_code"`
	ParentsContact int64  `json:"parents_contact"`
	FullAddress    string `json:"full_address"`
}

type StudentUserResponse struct {
	StudentID      *string `json:"student_id,omitempty"`
	FullName       *string `json:"full_name,omitempty"`
	Email          *string `json:"email,omitempty"`
	CreatedAt      *int    `json:"created_at,omitempty"`
	Password       *string `json:"-"`
	ContactNumber  *int64  `json:"contact_number,omitempty"`
	Class          *string `json:"class,omitempty"`
	TeacherCode    *string `json:"-"`
	ParentsContact *int64  `json:"parents_contact,omitempty"`
	FullAddress    *string `json:"full_address,omitempty"`
}

type StudentModelAddGroup struct {
	OwnerName      string `json:"owner_name"`
	GroupID        string `json:"group_id"`
	StudentID      string `json:"student_id"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	CreatedAt      int    `json:"created_at"`
	ContactNumber  int64  `json:"contact_number"`
	Class          string `json:"class"`
	ParentsContact int64  `json:"parents_contact"`
	FullAddress    string `json:"full_address"`
}
