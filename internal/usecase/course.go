package usecase

import (
	"context"
	"fmt"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/domain/models"
	"lms-bootcamp/internal/repository"

	"github.com/gosimple/slug"
)

type CourseUseCase struct {
	repo *repository.CourseRepository
}


func NewCourseUseCase(repo *repository.CourseRepository) *CourseUseCase {
	return &CourseUseCase{
		repo: repo,
	}
}

func (uc *CourseUseCase) CreateCourse(ctx context.Context, course *dto.CourseRequest) error {
	mentor, errMentor := uc.repo.ParallelFindUser(ctx, course.Mentor, dto.RoleMentor)
	student, errStudent := uc.repo.ParallelFindUser(ctx, course.Student, dto.RoleStudent)

	if errMentor != nil || len(mentor) == 0 {
		return errMentor
	}

	if errStudent != nil || len(student) == 0 {
		return errStudent
	}
	
	err := uc.repo.Create(ctx, &models.CourseModel{
		Title:       course.Title,
		Description: course.Description,
		Slug:       slug.Make(course.Title),
		Mentors:     mentor,
		Students:    student,
	})

	if(err != nil) {
		return err
	}

	return nil
}

func (uc *CourseUseCase) GetCourseByIDOrSlug(ctx context.Context, idOrSlug string) (*models.CourseModel, error) {
	return uc.repo.FindDetail(ctx, idOrSlug)
}

func (uc *CourseUseCase) UpdateCourse(ctx context.Context, id string, course *dto.CourseRequest) error {
	return uc.repo.Update(ctx, id, &models.CourseModel{
		Title:       course.Title,
		Description: course.Description,
		Slug:       slug.Make(course.Title),
	})
}

func (uc *CourseUseCase) DeleteCourse(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *CourseUseCase) GetAllCourses(ctx context.Context, filter *dto.CourseFilter, page, pageSize *int) (*dto.PaginationResponse[models.CourseModel], error) {
	fmt.Println("Page:", *page)
	fmt.Println("PageSize:", *pageSize)
	_, err := uc.repo.GetByFilter(ctx, filter, page, pageSize, "courses")
	if err != nil {
		return nil, err
	}
	return uc.repo.GetAll(ctx, page, pageSize)
}

func (uc *CourseUseCase) GetCoursesByUser(ctx context.Context, userID string, role string) ([]*models.UserModel, error) {
	return uc.repo.GetAllUser(ctx, userID, &role)
}
