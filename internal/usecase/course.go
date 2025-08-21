package usecase

import (
	"context"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/domain/models"
	"lms-bootcamp/internal/repository"

	"github.com/gosimple/slug"
)

type CourseUseCase struct {
	repository *repository.CourseRepository
}


func NewCourseUseCase(repo *repository.CourseRepository) *CourseUseCase {
	return &CourseUseCase{
		repository: repo,
	}
}

func (uc *CourseUseCase) CreateCourse(ctx context.Context, course *dto.CourseRequest) error {
	mentor, errMentor := uc.repository.ParallelFindUser(ctx, course.Mentor, dto.RoleMentor)
	student, errStudent := uc.repository.ParallelFindUser(ctx, course.Student, dto.RoleStudent)

	if errMentor != nil || len(mentor) == 0 {
		return errMentor
	}

	if errStudent != nil {
		return errStudent
	}

	err := uc.repository.Add(ctx, &models.CourseModel{
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
	return uc.repository.FindDetail(ctx, idOrSlug)
}

func (uc *CourseUseCase) UpdateCourse(ctx context.Context, id string, course *dto.CourseRequest) error {
	return uc.repository.Update(ctx, id, &models.CourseModel{
		Title:       course.Title,
		Description: course.Description,
		Slug:       slug.Make(course.Title),
	})
}

func (uc *CourseUseCase) DeleteCourse(ctx context.Context, id string) error {

	return uc.repository.Delete(ctx, id)
}

func (uc *CourseUseCase) GetAllCourses(ctx context.Context, search string, page, pageSize *int) (*dto.PaginationResponse[models.CourseModel], error) {
	_, err := uc.repository.GetByFilter(ctx, search, page, pageSize)
	if err != nil {
		return nil, err
	}
	return uc.repository.GetAll(ctx, page, pageSize)
}

func (uc *CourseUseCase) GetCoursesByUser(ctx context.Context, userID string, role string) ([]*dto.UserData, error) {
	return uc.repository.GetAllUser(ctx, userID, &role)
}

func (uc *CourseUseCase) BulkAddUser(ctx context.Context, courseId string, role *string, userIds []*string) error {
	return uc.repository.BulkAddUser(ctx, courseId, role, userIds)
}

func (uc *CourseUseCase) RemoveUser(ctx context.Context, courseId string, userId *string, role string) error {
	return uc.repository.RemoveUser(ctx, courseId, *userId, &role)
}
