package usecase

import (
	"context"
	"lms-bootcamp/internal/domain/dto"
)

type RoleUseCase interface {
	CreateRole(ctx context.Context, name string) (string, error)
	GetRoleByID(ctx context.Context, id string) (*dto.RoleData, error)
	GetAllRoles(ctx context.Context) ([]*dto.RoleData, error)
	UpdateRole(ctx context.Context, id string, name string) error
	DeleteRole(ctx context.Context, id string) error
}