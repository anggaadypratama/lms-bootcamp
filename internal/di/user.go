package di

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/repository"
	"lms-bootcamp/internal/usecase"

	"gorm.io/gorm"
)

type UserDI struct {
	db *gorm.DB
	Repo    *repository.UserRepository
	Usecase *usecase.UserUseCase
	Handler *handler.UserHandler
}

func NewUserDI(db *gorm.DB) *UserDI {
	repo := repository.NewUserRepository(db)
	usecase := usecase.NewUserUseCase(repo)
	handler := handler.NewUserHandler(*usecase)

	return &UserDI{
		db:      db,
		Repo:    repo,
		Usecase: usecase,
		Handler: handler,
	}
}

func (di *UserDI) Init() *handler.UserHandler {
	return di.Handler
}
