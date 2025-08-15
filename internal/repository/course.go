package repository

import (
	"context"
	"fmt"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/domain/models"
	"lms-bootcamp/internal/pkg/utils"
	"reflect"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository struct {
	base *GormRepository[models.CourseModel]
	user *GormRepository[models.UserModel]
    utils *utils.Utils
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
    utils := utils.NewUtils()
	return &CourseRepository{
		base: NewGormRepository[models.CourseModel](db),
		user: NewGormRepository[models.UserModel](db),
		utils: utils,
	}
}

func (r *CourseRepository) FindUser(ctx context.Context, userId string, roleName string) ([]models.UserModel, error) {
	var users []models.UserModel
	err := r.user.db.Joins("JOIN roles ON roles.id = users.role_id").Where("users.id = ? AND roles.name = ? AND users.deleted_at IS NULL AND roles.deleted_at IS NULL", userId, roleName).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *CourseRepository) ParallelFindUser(ctx context.Context, ids []string, role string) ([]*models.UserModel, error) {
    var result []*models.UserModel
    var mu sync.Mutex
    var wg sync.WaitGroup
    var errUser error

    for _, id := range ids {
        wg.Add(1)
        go func(userID string) {
            defer wg.Done()
            res, err := r.FindUser(ctx, userID, role)
            mu.Lock()
            defer mu.Unlock()
            if err != nil || len(res) == 0 {
                if errUser == nil {
                    if err != nil {
                        errUser = err
                    } else {
                        errUser = fmt.Errorf("user %s dengan role %s tidak ditemukan", userID, role)
                    }
                }
                return
            }
            result = append(result, &res[0])
        }(id)
    }
    wg.Wait()
    if errUser != nil {
        return nil, errUser
    }
    return result, nil
}

func (r *CourseRepository) Add(ctx context.Context, course *models.CourseModel) error {
    return r.base.Create(ctx, course)
}

func (r *CourseRepository) FindDetail(ctx context.Context, idOrSlug string) (*models.CourseModel, error) {
    var course models.CourseModel
    db := r.base.db
    if _, err := uuid.Parse(idOrSlug); err == nil {
        err = db.Where("id = ? AND deleted_at IS NULL", idOrSlug).First(&course).Error
        if err != nil {
            return nil, err
        }
    } else {
        err = db.Where("slug = ? AND deleted_at IS NULL", idOrSlug).First(&course).Error
        if err != nil {
            return nil, err
        }
    }
    return &course, nil
}

func (r *CourseRepository) Update(ctx context.Context, id string, course *models.CourseModel) error {
    return r.base.Update(ctx, id, course)
}

func (r *CourseRepository) Delete(ctx context.Context, id string) error {
    return r.base.Delete(ctx, id)
}

func (r *CourseRepository) GetAll(ctx context.Context, page, pageSize *int) (*dto.PaginationResponse[models.CourseModel], error) {
	if page == nil && pageSize == nil {
		courses, err := r.base.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		return dto.NewPaginationResponse(courses, int64(len(courses)), 1, len(courses)), nil
	}

	courses, err := r.base.Paginate(ctx, *page, *pageSize)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *CourseRepository) GetByFilter(ctx context.Context, filter *dto.CourseFilter, page, pageSize *int, table string) (*dto.PaginationResponse[models.CourseModel], error) {
    session := r.base.db.Session(&gorm.Session{})
    sessionRepo := NewGormRepository[models.CourseModel](session)

    if filter.Title != "" {
        sessionRepo = sessionRepo.Where(table+".title ILIKE ?", "%"+filter.Title+"%")
    }

    if page == nil || pageSize == nil {
        courses, err := sessionRepo.GetAll(ctx)
        if err != nil {
            return nil, err
        }
        return dto.NewPaginationResponse(courses, int64(len(courses)), 1, len(courses)), nil
    }

    courses, err := sessionRepo.Paginate(ctx, *page, *pageSize)
    if err != nil {
        return nil, err
    }

    fmt.Println("courses:", *courses)
    return courses, nil
}

func (r *CourseRepository) GetAllUser(ctx context.Context, id string, role *string) ([]*dto.UserData, error) {
    var course models.CourseModel

    err := r.base.db.Preload("Students").Preload("Mentors").First(&course, "id = ?", id).Where("courses.deleted_at IS NULL").Error

    v := reflect.ValueOf(course)
    field := v.FieldByName(r.utils.ToPascalCase(*role)+"s")
    if err != nil {
        return nil, err
    }

    if !field.IsValid() {
        return nil, fmt.Errorf("field %s not found in course model", *role)
    }
    

    users, ok := field.Interface().([]*models.UserModel)
    if !ok {
        return nil, fmt.Errorf("field %s is not a slice of UserModel", *role)
    }

    var userData []*dto.UserData
    for _, user := range users {
        userData = append(userData, &dto.UserData{
            ID:   user.ID,
            Name: user.Name,
            Email: user.Email,
            PhoneNumber: user.PhoneNumber,
        })
    }

    return userData, nil
}

func (r *CourseRepository) BulkAddUser(ctx context.Context, courseId string, role *string, userIds []*string) error {
    return r.base.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var course models.CourseModel

        if err := r.base.db.First(&course, "id = ? AND deleted_at IS NULL", courseId).Error; err != nil {
            return fmt.Errorf("course with ID %s not found", courseId)
        }

        var users []*models.UserModel
        if err := r.user.db.WithContext(ctx).Joins("Role").Where("users.deleted_at IS NULL AND users.id IN ?", userIds).Find(&users).Error; err != nil {
            return fmt.Errorf("failed to find users: %v", err)
        }

        for _, user := range users {
            if user.Role.Name != *role {
                return fmt.Errorf("user %s is not a %s", user.ID, *role)
            }
        }

        if err := r.base.db.Model(&course).Association(r.utils.ToPascalCase(*role) + "s").Append(users); err != nil {
            return fmt.Errorf("failed to add mentors: %v", err)
        }

        return nil
	})

}

func (r *CourseRepository) RemoveUser(ctx context.Context, courseId string, userId string, role *string) error {
    return r.base.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var course models.CourseModel
        if err := tx.First(&course, "id = ? AND deleted_at IS NULL", courseId).Error; err != nil {
            return fmt.Errorf("course with ID %s not found", courseId)
        }

        var user models.UserModel
        if err := tx.First(&user, "id = ? AND role_id = (SELECT id FROM roles WHERE name = ?) AND deleted_at IS NULL", userId, *role).Error; err != nil {
            return fmt.Errorf("user with ID %s and role %s not found", userId, *role)
        }

        if err := tx.Model(&course).Association(r.utils.ToPascalCase(*role) + "s").Delete(&user); err != nil {
            return fmt.Errorf("failed to remove user: %v", err)
        }

        return nil
    })
}