package noticeboard

type NoticeBoardModel struct {
	NoticeID    *string `json:"notice_id,omitempty"`
	Title       *string `json:"title,omitempty"`
	Desctiption *string `json:"description,omitempty"`
	GroupID     string `json:"-"`
	TeacherID   string `json:"-"`
	CreatedAt   int64  `json:"-"`
	UpdatedAt   *int64  `json:"updated_at,omitempty"`
}
