package repository

import (
	"context"
	"fmt"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/domain/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	base *GormRepository[models.UserModel]
	role *GormRepository[models.RoleModel]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		base: NewGormRepository[models.UserModel](db),
		role: NewGormRepository[models.RoleModel](db),
	}
}

func (r *UserRepository) GetRole(ctx context.Context, id string) (*models.RoleModel, error) {
	return r.role.GetById(ctx, id, "roles")
}

func (r *UserRepository) Create(ctx context.Context, user *models.UserModel) error {
    return r.base.Create(ctx, user)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.UserModel, error) {
	return r.base.Join("Role").GetById(ctx, id, "users")
}

func (r *UserRepository) Update(ctx context.Context, id string, user *models.UserModel) error {
    return r.base.Update(ctx, id, user)
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
    return r.base.Delete(ctx, id)
}

func (r *UserRepository) GetAll(ctx context.Context, page, pageSize *int) (*dto.PaginationResponse[models.UserModel], error) {
	if page == nil && pageSize == nil {
		users, err := r.base.Where("users.deleted_at IS NULL").Join("Role").GetAll(ctx)
		if err != nil {
			return nil, err
		}
		return dto.NewPaginationResponse(users, int64(len(users)), 1, len(users)), nil
	}

	users, err := r.base.Where("users.deleted_at IS NULL").Join("Role").Paginate(ctx, *page, *pageSize)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetByFilter(ctx context.Context, filter *map[string]interface{}, page, pageSize *int, table string) (*dto.PaginationResponse[models.UserModel], error) {
	defaultDb := r.base.Join("Role").Where("users.deleted_at IS NULL")

    if filter != nil {
        for key, value := range *filter {
            if key == "name" || key == "email" {
                defaultDb = defaultDb.Where(table+"."+key+" ILIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
            } else {
                defaultDb = defaultDb.Where(table+"."+key+" = ?", value)
            }
        }
    }

	if page == nil && pageSize == nil  {
		users, err := defaultDb.GetAll(ctx, "users")
		if err != nil {
			return nil, err
		}
		return dto.NewPaginationResponse(users, int64(len(users)), 1, len(users)), nil
	}

	users, err := defaultDb.Paginate(ctx, *page, *pageSize)
	if err != nil {
		return nil, err
	}
	return users, nil
}