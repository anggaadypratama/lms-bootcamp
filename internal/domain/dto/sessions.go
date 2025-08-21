package dto

type SessionRequest struct {
	CourseID 	string `json:"course_id" binding:"required" jsonschema:"required,description=Course ID"`
	Title	 	string `json:"title" binding:"required" jsonschema:"required,description=Session title"`
	Description string `json:"description" binding:"required" jsonschema:"required,description=Session description"`
	StartTime   string `json:"start_time" binding:"required" jsonschema:"required,description=Session start time"`
	EndTime     string `json:"end_time" binding:"required" jsonschema:"required,description=Session end time"`
}

type SessionData struct {
	ID          string `json:"id"`
	CourseID   	string `json:"course_id"`
	Title      	string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Course      *SessionCourse `json:"course,omitempty"`
}

type SessionCourse struct {
	Title string `json:"title"`
	Mentor []*SessionMentor `json:"mentor,omitempty"`
}

type SessionMentor struct {
	Name  string `json:"name"`
}