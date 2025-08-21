package usecase

import (
	"context"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/domain/models"
	"lms-bootcamp/internal/repository"
)


type RoleUseCase struct {
	repository *repository.RoleRepository
}

func NewRoleUseCase(repository *repository.RoleRepository) *RoleUseCase {
	return &RoleUseCase{
		repository: repository,
	}
}

func (uc *RoleUseCase) CreateRole(ctx context.Context, name string) (string, error) {
	err := uc.repository.Create(ctx, &models.RoleModel{Name: name})
	if err != nil {
		return "", err
	}
	return "Role created successfully", nil
}


func (uc *RoleUseCase) GetRoleByID(ctx context.Context, id string) (*dto.RoleData, error) {
	role, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.RoleData{
		ID:   role.ID,
		Name: role.Name,
	}, nil
}

func (uc *RoleUseCase) GetAllRoles(ctx context.Context) ([]*dto.RoleData, error) {
	roles, err := uc.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var roleData []*dto.RoleData
	for _, role := range roles {
		roleData = append(roleData, &dto.RoleData{
			ID:   role.ID,
			Name: role.Name,
		})
	}
	return roleData, nil
}

func (uc *RoleUseCase) UpdateRole(ctx context.Context, id string, name string) error {
	role, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	role.Name = name
	err = uc.repository.Update(ctx, id, role)
	if err != nil {
		return err
	}
	return nil
}


func (uc *RoleUseCase) DeleteRole(ctx context.Context, id string) error {
	_, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	return uc.repository.Delete(ctx, id)
}