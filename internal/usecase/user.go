package usecase

import (
	"context"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/domain/models"
	"lms-bootcamp/internal/repository"

	"golang.org/x/crypto/bcrypt"
)


type UserUseCase struct {
	repository *repository.UserRepository
}

func NewUserUseCase(repository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repository: repository,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, item *dto.UserRequest) (*string, error) {
	role, err := uc.repository.GetRole(ctx, item.RoleId)
	if err != nil {
		return nil, err
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return nil, hashErr
	}
	err = uc.repository.Create(ctx, &models.UserModel{
		Name:        item.Name,
		Email:       item.Email,
		Password:    string(hashedPassword),
		PhoneNumber: item.PhoneNumber,
		RoleId:      role.ID,
	})

	if err != nil {
		return nil, err
	}
	msg := "User created successfully"
	return &msg, nil
}


func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (*dto.UserData, error) {
	user, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.UserData{
			ID:   user.ID,
			Name: user.Name,
			Email: user.Email,
			PhoneNumber: user.PhoneNumber,
			Role: &dto.RoleData{
				ID:   user.RoleId,
				Name: user.Role.Name,
			},
		}, nil
}

func (uc *UserUseCase) GetAllUsers(ctx context.Context, filter *map[string]interface{}, page, pageSize *int) (*dto.PaginationResponse[*dto.UserData], error) {
	usersPage, err := uc.repository.GetByFilter(ctx, filter, page, pageSize, "users")
	if err != nil {
		return nil, err
	}

	var userData []*dto.UserData
	for _, user := range usersPage.Data {
		userData = append(userData, &dto.UserData{
			ID:   user.ID,
			Name: user.Name,
			Email: user.Email,
			PhoneNumber: user.PhoneNumber,
			Role: &dto.RoleData{
				ID:   user.RoleId,
				Name: user.Role.Name,
			},
		})
	}
	return dto.NewPaginationResponse(userData, usersPage.Total, usersPage.Page, usersPage.PerPage), nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, id string,  item *dto.UserRequest) error {
	user, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user.Role == nil {
		return err
	}

	password := item.Password
	if password != "" {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		password = string(hashedPassword)
	}

	user = &models.UserModel{
		Name:      item.Name,
		Email:    item.Email,
		Password:  password,
		PhoneNumber: item.PhoneNumber,
		RoleId:  user.RoleId,
	}

	err = uc.repository.Update(ctx, id, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	err := uc.repository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}