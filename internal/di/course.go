package di

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/repository"
	"lms-bootcamp/internal/usecase"

	"gorm.io/gorm"
)

type CourseDI struct {
	db *gorm.DB
	Repo    *repository.CourseRepository
	Usecase *usecase.CourseUseCase
	Handler *handler.CourseHandler
}

func NewCourseDI(db *gorm.DB) *CourseDI {
	repo := repository.NewCourseRepository(db)
	usecase := usecase.NewCourseUseCase(repo)
	handler := handler.NewCourseHandler(usecase)

	return &CourseDI{
		db:      db,
		Repo:    repo,
		Usecase: usecase,
		Handler: handler,
	}
}

func (di *CourseDI) Init() *handler.CourseHandler {
	return di.Handler
}
