package notes

type SaveNotes struct {
	NotesTitle       string `json:"notes_title,omitempty"`
	NotesDescription string `json:"notes_desc,omitempty"`
	ClassName        string `json:"class_name,omitempty"`
	NotesTopic       string `json:"notes_topic,omitempty"`
	NotesVisiblity   string `json:"notes_visiblity,omitempty"`
	IsEditable       bool   `json:"is_editable,omitempty"`
	GroupId          string `json:"group_id,omitempty"`
}

type GroupNotes struct {
	NotesId          string   `json:"notes_id,omitempty"`
	NotesTitle       string   `json:"notes_title,omitempty"`
	NotesDescription string   `json:"notes_desctription,omitempty"`
	ClassName        string   `json:"class_name,omitempty"`
	NotesTopic       string   `json:"notes_topic,omitempty"`
	NotesSubject     string   `json:"notes_subject,omitempty"`
	TeacherId        string   `json:"teacher_id,omitempty"`
	UploadedAt       int64    `json:"uploaded_at,omitempty"`
	FileURLsList     []string `json:"attached_files,omitempty"`
	FileNames        string   `json:"-"`
	NotesVisiblity   string   `json:"notes_visiblity,omitempty"`
	IsEditableInt    int      `json:"-"`
	IsEditable       bool     `json:"is_editable,omitempty"`
	GroupId          string   `json:"group_id,omitempty"`
}
