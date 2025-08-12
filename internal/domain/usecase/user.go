package usecase

import (
	"context"
	"lms-bootcamp/internal/domain/dto"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, item *dto.UserRequest) (*string, error)
	GetUserByID(ctx context.Context, id string) (*dto.UserData, error)
	GetAllUsers(ctx context.Context, filter *map[string]interface{}, page, pageSize *int) (*dto.PaginationResponse[*dto.UserData], error)
	UpdateUser(ctx context.Context, id string, item *dto.UserRequest) error
	DeleteUser(ctx context.Context, id string) error
}