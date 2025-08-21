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

func (r *CourseRepository) GetByFilter(ctx context.Context, search string, page, pageSize *int) (*dto.PaginationResponse[models.CourseModel], error) {
    sessionRepo := r.base.Where("courses.deleted_at IS NULL")

    if search != "" {
        sessionRepo = sessionRepo.Where("title ILIKE ?", "%"+search+"%")
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
    return courses, nil
}

func (r *CourseRepository) GetAllUser(ctx context.Context, id string, role *string) ([]*dto.UserData, error) {
    var course models.CourseModel
    var roles *models.RoleModel

    if err := r.base.db.First(&roles, "id = ?", *role).Error; err != nil {
        return nil, fmt.Errorf("role %s not found", *role)
    }
    
    err := r.base.db.Preload("Students").Preload("Mentors").First(&course, "id = ?", id).Where("courses.deleted_at IS NULL").Error

    v := reflect.ValueOf(course)
    field := v.FieldByName(r.utils.ToPascalCase(roles.Name)+"s")
    if err != nil {
        return nil, err
    }

    if !field.IsValid() {
        return nil, fmt.Errorf("field %s not found in course model", roles.Name)
    }
    

    users, ok := field.Interface().([]*models.UserModel)
    if !ok {
        return nil, fmt.Errorf("field %s is not a slice of UserModel", roles.Name)
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

func (r *CourseRepository) BulkAddUser(ctx context.Context, courseId string, role_id *string, userIds []*string) error {
    return r.base.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var course *models.CourseModel
        var role *models.RoleModel
        
        if err := r.base.db.First(&role, "roles.id = ?", *role_id).Error; err != nil {
            return fmt.Errorf("role %s not found", *role_id)
        }

        if err := r.base.db.First(&course, "courses.id = ? AND courses.deleted_at IS NULL", courseId).Error; err != nil {
            return fmt.Errorf("course with ID %s not found", courseId)
        }

        var users []*models.UserModel
            if err := r.user.db.WithContext(ctx).Preload("Role").Where("users.role_id = ? AND users.id IN ?", *role_id, userIds).Find(&users).Error; err != nil {
                return fmt.Errorf("failed to find user: %v", err)
            }

            // Restore soft deleted users (deleted_at != null)
            for _, user := range users {
                if user.DeletedAt.Valid {
                    if err := r.user.db.Model(user).Update("deleted_at", nil).Error; err != nil {
                        return fmt.Errorf("failed to restore user %s: %v", user.ID, err)
                    }
                    user.DeletedAt.Valid = false
                }
            }

        for _, user := range users {
            if user.Role.ID != *role_id {
                return fmt.Errorf("user %s is not a %s", user.ID, role.Name)
            }
        }

        if err := r.base.db.Model(&course).Association(r.utils.ToPascalCase(role.Name) + "s").Append(users); err != nil {
            return fmt.Errorf("failed to add user: %v", err)
        }

        return nil
	})

}

func (r *CourseRepository) RemoveUser(ctx context.Context, courseId string, userId string, role *string) error {
    return r.base.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var course models.CourseModel
        var roles *models.RoleModel

        if err := r.base.db.First(&roles, "id = ?", *role).Error; err != nil {
            return fmt.Errorf("role %s not found", *role)
        }

        if err := tx.First(&course, "id = ? AND deleted_at IS NULL", courseId).Error; err != nil {
            return fmt.Errorf("course with ID %s not found", courseId)
        }

        var user models.UserModel
        if err := tx.First(&user, "id = ? AND role_id = (SELECT id FROM roles WHERE name = ?) AND deleted_at IS NULL", userId, roles.Name).Error; err != nil {
            return fmt.Errorf("user with ID %s and role %s not found", userId, roles.Name)
        }

        if err := tx.Model(&course).Association(r.utils.ToPascalCase(roles.Name) + "s").Delete(&user); err != nil {
            return fmt.Errorf("failed to remove user: %v", err)
        }

        return nil
    })
}