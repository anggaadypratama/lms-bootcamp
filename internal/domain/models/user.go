package models

type UserModel struct {
	BaseModel
	Email        string     `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Password     string     `json:"password" gorm:"not null;size:255"`
	PhoneNumber  string     `json:"phone_number" gorm:"size:20"`
	Name         string     `json:"name" gorm:"not null;size:255"`
	RoleId      string     `json:"role_id" gorm:"type:varchar(36);not null;index;"`
	Role        *RoleModel `json:"role,omitempty" gorm:"foreignKey:RoleId;references:ID"`
	EnrolledCourses []*CourseModel `json:"enrolled_courses,omitempty" gorm:"many2many:course_users;"`
	MentorCourses   []*CourseModel `json:"mentor_courses,omitempty" gorm:"many2many:course_mentors;"`
}

func (UserModel) TableName() string {
	return "users"
}