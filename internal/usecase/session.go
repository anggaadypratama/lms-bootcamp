package usecase

import (
	"context"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/domain/models"
	"lms-bootcamp/internal/repository"
)


type SessionUseCase struct {
	repository *repository.SessionRepository
}

func NewSessionUseCase(repository *repository.SessionRepository) *SessionUseCase {
	return &SessionUseCase{
		repository: repository,
	}
}

func (uc *SessionUseCase) CreateSession(ctx context.Context, item *dto.SessionRequest) (string, error) {
	err := uc.repository.Create(ctx, &models.SessionModel{
		CourseId:   &item.CourseID,
		Title:      item.Title,
		Description: item.Description,
		StartTime:   &item.StartTime,
		EndTime:     &item.EndTime,
	})
	if err != nil {
		return "", err
	}
	return "Session created successfully", nil
}


func (uc *SessionUseCase) GetSessionByID(ctx context.Context, id string) (*dto.SessionData, error) {
	session, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	course, err := uc.repository.GetMentor(ctx, *session.CourseId)
	if err != nil {
		return nil, err
	}

	return &dto.SessionData{
		ID: session.ID,
		Title:      session.Title,
		CourseID:  *session.CourseId,
		Description: session.Description,
		StartTime:  *session.StartTime,
		EndTime:    *session.EndTime,
		Course: &dto.SessionCourse{
			Title: session.Course.Title,
			Mentor: func() []*dto.SessionMentor {
				var mentors []*dto.SessionMentor
				for _, mentor := range course.Mentors{
					mentors = append(mentors, &dto.SessionMentor{Name: mentor.Name})
				}
				return mentors
			}(),
		},
	}, nil
}

func (uc *SessionUseCase) GetAllSessions(ctx context.Context, search string, page, pageSize int) (*dto.PaginationResponse[*dto.SessionData], error) {

	total, err := uc.repository.CountTotalSessions(ctx, search)
	if err != nil {
		return nil, err
	}
	
	if pageSize == 0 {
		pageSize = int(*total)
	}

	if page == 0 {
		page = 1
	}

	sessions, err := uc.repository.GetAll(ctx, search, page, pageSize)
	if err != nil {
		return nil, err
	}



	var sessionData []*dto.SessionData
	for _, session := range sessions {
		sessionData = append(sessionData, &dto.SessionData{
			ID: session.ID,
			Title:      session.Title,
			Description: session.Description,
			StartTime: *session.StartTime,
			EndTime:   *session.EndTime,
			Course:  &dto.SessionCourse{
			Title: session.Course.Title,
			Mentor: func() []*dto.SessionMentor {
				var mentors []*dto.SessionMentor
				for _, mentor := range session.Course.Mentors {
					mentors = append(mentors, &dto.SessionMentor{Name: mentor.Name})
				}
				return mentors
			}(),
		},
	})
	}
	return dto.NewPaginationResponse(sessionData, *total, page, pageSize), nil
}

func (uc *SessionUseCase) UpdateSession(ctx context.Context, id string,  item *dto.SessionRequest) error {
	session, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	session.CourseId = &item.CourseID
	session.Title = item.Title
	session.Description = item.Description
	session.StartTime = &item.StartTime
	session.EndTime = &item.EndTime

	err = uc.repository.Update(ctx, id, session)
	if err != nil {
		return err
	}
	return nil
}


func (uc *SessionUseCase) DeleteSession(ctx context.Context, id string) error {
	_, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.repository.Delete(ctx, id)
}