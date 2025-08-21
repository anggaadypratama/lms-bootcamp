package dto

type CourseRequest struct {
		Title       string   `json:"title" binding:"required" jsonschema:"required,description=Course title"`
		Description string   `json:"description" binding:"required" jsonschema:"required,description=Course description"`
		Mentor      []string `json:"mentor" binding:"required" jsonschema:"required,description=Course mentor"`
		Student     []string `json:"student" binding:"required" jsonschema:"required,description=Course student"`
}

type CourseUserRequest struct {
	User     []*string `json:"ids" binding:"required" jsonschema:"required,description=Course user"`
	RoleID  *string    `json:"role_id" binding:"required" jsonschema:"required,description=Course role ID"`
}

type CourseQuery struct {
	Role []string `form:"role"`
}

type CourseFilter struct {
	Title string `form:"title"`
}