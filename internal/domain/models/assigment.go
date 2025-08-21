package models

type AssignmentModel struct {
	BaseModel
	StudentId 	 string `json:"student_id,omitempty" gorm:"foreignKey:StudentId;references:ID"`
	Student      *UserModel `json:"student,omitempty" gorm:"foreignKey:StudentId;references:ID"`
	FileUrl		string    `json:"file_url" gorm:"not null;size:255"`
	SubmittedAt *string   `json:"submitted_at,omitempty" gorm:"type:timestamp"`
	Status 		string    `json:"status" gorm:"not null;default:'pending';size:20"`
	SessionId  	string     `json:"session_id" gorm:"not null;index"`
	Grade 		string	`json:"grade,omitempty" gorm:"size:20"`
	Session   *SessionModel `json:"session,omitempty" gorm:"foreignKey:SessionId;references:ID"`
}

func (AssignmentModel) TableName() string {
	return "assignments"
}