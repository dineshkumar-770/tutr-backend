package creategroup

type CreateGroupModel struct {
	GroupID          *string                   `json:"group_id,omitempty"`
	TeacherID        *string                   `json:"teacher_id,omitempty"`
	CreatedAt        *int64                    `json:"created_at,omitempty"`
	GroupClass       *string                   `json:"group_class,omitempty"`
	GroupName        *string                   `json:"group_name,omitempty"`
	GroupDescription *string                   `json:"group_desc,omitempty"`
	AllMembers       *[]GroupMemberStudentsData `json:"all_members,omitempty"`
}
