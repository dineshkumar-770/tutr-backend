package creategroup

import "main/models"

type AllTeacherByStudentModel struct {
	GroupId          *string                     `json:"group_id,omitempty"`
	GroupName        *string                     `json:"group_name,omitempty"`
	GroupDesctiption *string                     `json:"group_description,omitempty"`
	TeacherId        *string                     `json:"teacher_id,omitempty"`
	CreatedAt        *int64                      `json:"created_at,omitempty"`
	GroupClass       *string                     `json:"group_class,omitempty"`
	TeacherDetails   *models.Teacher             `json:"teacher_details,omitempty"`
	// StudentDetails   *models.StudentUserResponse `json:"student_details,omitempty"`
}
