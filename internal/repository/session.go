package repository

import (
	"context"
	"lms-bootcamp/internal/domain/models"

	"gorm.io/gorm"
)

type SessionRepository struct {
	base *GormRepository[models.SessionModel]
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		base: NewGormRepository[models.SessionModel](db),
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.SessionModel) error {
    return r.base.Create(ctx, session)
}

func (r *SessionRepository) FindByID(ctx context.Context, id string) (*models.SessionModel, error) {
    var session *models.SessionModel
    err := r.base.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        err := r.base.db.Joins("Course").Where("sessions.id = ?", id).First(&session).Error
        if err != nil {
            return err
        }

        return nil
    })
    if err != nil {
        return nil, err
    }
    return session, nil
}

func (r *SessionRepository) GetMentor(ctx context.Context, course_id string) (*models.CourseModel, error) {
    var course *models.CourseModel
    err := r.base.db.WithContext(ctx).
        Preload("Mentors").
        Where("id = ?", course_id).
        First(&course).Error
    if err != nil {
        return nil, err
    }
    return course, nil
}

func (r *SessionRepository) Update(ctx context.Context, id string, session *models.SessionModel) error {
    return r.base.Update(ctx, id, session)
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
    return r.base.Delete(ctx, id)
}

func (r *SessionRepository) GetAll(ctx context.Context, search string, page, pageSize int) ([]models.SessionModel, error) {
    var sessions []models.SessionModel
    db := r.base.db.WithContext(ctx).Joins("Course").Where("sessions.deleted_at IS NULL")


    if search != "" {
        db = db.Where("sessions.title ILIKE ?", "%"+search+"%")
    }

    if page != 0 && pageSize != 0 {
        db = db.Offset((page-1)*(pageSize)).Limit(pageSize)
    }

    err := db.Find(&sessions).Error
    if err != nil {
        return nil, err
    }
    return sessions, nil
}

func (r *SessionRepository) CountTotalSessions(ctx context.Context, search string) (*int64, error) {
    var total int64

    countDB := r.base.db.Session(&gorm.Session{})

    if search != "" {
        countDB = countDB.Where("sessions.title ILIKE ?", "%"+search+"%")
    }
    if err := countDB.WithContext(ctx).Model(new(models.SessionModel)).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
        return nil, err
    }
    

    return &total, nil

}