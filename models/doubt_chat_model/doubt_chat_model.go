package doubtchatmodel

type DoubtChatDataModel struct {
	DoubtID     string `json:"doubt_id,omitempty"`
	GroupID     string `json:"group_id,omitempty"`
	StudentID   string `json:"student_id,omitempty"`
	TeacherID   string `json:"teacher_id,omitempty"`
	DoubtText   string `json:"doubt_text,omitempty"`
	FilePath    string `json:"file_url,omitempty"`
	DoubtStatus string `json:"status,omitempty"`
	CreatedAt   int    `json:"created_at,omitempty"`
	UpdatedAt   int    `json:"updated_at,omitempty"`
}

type DoubtChatData struct {
	DoubtID         *string `json:"doubt_id,omitempty"`
	GroupID         *string `json:"group_id,omitempty"`
	StudentID       *string `json:"student_id,omitempty"`
	TeacherID       *string `json:"teacher_id,omitempty"`
	DoubtText       *string `json:"doubt_text,omitempty"`
	FilePath        *string `json:"file_path,omitempty"`
	DoubtStatus     *string `json:"status,omitempty"`
	CreatedAt       *int64  `json:"created_at,omitempty"`
	UpdatedAt       *int64  `json:"updated_at,omitempty"`
	StudentFullName *string `json:"full_name,omitempty"`
	StudentEmail    *string `json:"student_email,omitempty"`
	StudentClass    *string `json:"student_class,omitempty"`
}
