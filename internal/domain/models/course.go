package models

type CourseModel struct {
	BaseModel
	Title       string `json:"title" gorm:"not null;size:255"`
	Description string `json:"description" gorm:"type:text"`
    Students    []*UserModel `json:"students,omitempty" gorm:"many2many:course_users;"`
    Mentors     []*UserModel `json:"mentors,omitempty" gorm:"many2many:course_mentors;"`
}

func (CourseModel) TableName() string {
	return "courses"
}
