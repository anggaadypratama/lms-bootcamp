package dto

type CourseRequest struct {
		Title       string   `json:"title"`
		Description string   `json:"desc"`
		Mentor      []string `json:"mentor"`
		Student     []string `json:"student"`
}

type CourseFilter struct {
	Title string `form:"title"`
}