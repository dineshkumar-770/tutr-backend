package models

type Teacher struct {
	TeacherID       string `json:"teacher_id"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	ContactNumber   int64  `json:"contact_number"` // Using int64 for contact numbers
	Subject         string `json:"subject"`
	CreatedAt       int64  `json:"created_at"` // Timestamp
	Qualification   string `json:"qualification"`
	ExperienceYears int    `json:"experience_years"`
	Address         string `json:"address"`
	ClassAssigned   string `json:"class_assigned"`
	TeacherCode     string `json:"teacher_code"`
}
