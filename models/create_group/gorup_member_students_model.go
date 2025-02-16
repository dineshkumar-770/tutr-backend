package creategroup

type GroupMemberStudentsData struct {
	GroupMemberID           string `json:"group_member_id"`
	GroupID                 string `json:"group_id"`
	StudentID               string `json:"student_id"`
	OwnerID                 string `json:"owner_id"`
	GroupOwner              string `json:"group_owner_name"`
	StudentJoinedAt         int64  `json:"student_joined_at"`
	StudentFullName         string `json:"student_full_name"`
	StudentEmail            string `json:"student_email"`
	StudentAccountCreatedAt int64  `json:"student_account_creation_date"`
	StudentContactNumber    int64  `json:"student_contact"`
	StudentClass            string `json:"student_class"`
	StudentParentsContact   int64  `json:"student_parnets_contact"`
	StudentFullAddress      string `json:"student_full_address"`
}
