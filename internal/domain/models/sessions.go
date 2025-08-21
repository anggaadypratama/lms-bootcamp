package models

type SessionModel struct {
	BaseModel
	CourseId  *string `json:"course_id" gorm:"not null;index"`
	Title     string  `json:"title" gorm:"not null;size:255"`
	Description string  `json:"description" gorm:"not null"`
	StartTime *string  `json:"start_time,omitempty" gorm:"type:timestamp"`
	EndTime   *string  `json:"end_time,omitempty" gorm:"type:timestamp"`
	Course    *CourseModel `json:"course,omitempty" gorm:"foreignKey:CourseId;references:ID"`
	Assignment  []*AssignmentModel `json:"assignments,omitempty" gorm:"foreignKey:SessionId;references:ID"`
}

func (SessionModel) TableName() string {
	return "sessions"
}
