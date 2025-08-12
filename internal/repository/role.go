package repository

import (
	"context"
	"lms-bootcamp/internal/domain/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	base *GormRepository[models.RoleModel]
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		base: NewGormRepository[models.RoleModel](db),
	}
}

func (r *RoleRepository) Create(ctx context.Context, role *models.RoleModel) error {
    return r.base.Create(ctx, role)
}

func (r *RoleRepository) FindByID(ctx context.Context, id string) (*models.RoleModel, error) {
    return r.base.GetById(ctx, id)
}

func (r *RoleRepository) Update(ctx context.Context, id string, role *models.RoleModel) error {
    return r.base.Update(ctx, id, role)
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
    return r.base.Delete(ctx, id)
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]models.RoleModel, error) {
	return r.base.GetAll(ctx)
}