package di

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/pkg/utils"
	"lms-bootcamp/internal/repository"
	"lms-bootcamp/internal/usecase"

	"gorm.io/gorm"
)

type AuthDI struct {
	db *gorm.DB
	Repo    *repository.AuthRepository
	Usecase *usecase.AuthUseCase
	Handler *handler.AuthHandler
}

func NewAuthDI(db *gorm.DB) *AuthDI {
	repo := repository.NewAuthRepository(db)
	utilsInstance := utils.NewUtils()
	usecase := usecase.NewAuthUseCase(repo, utilsInstance)
	handler := handler.NewAuthHandler(usecase)

	return &AuthDI{
		db:      db,
		Repo:    repo,
		Usecase: usecase,
		Handler: handler,
	}
}

func (di *AuthDI) Init() *handler.AuthHandler {
	return di.Handler
}
