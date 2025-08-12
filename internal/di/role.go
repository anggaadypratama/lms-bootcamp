package di

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/repository"
	"lms-bootcamp/internal/usecase"

	"gorm.io/gorm"
)

type RoleDI struct {
	db *gorm.DB
	Repo    *repository.RoleRepository
	Usecase *usecase.RoleUseCase
	Handler *handler.RoleHandler
}

func NewRoleDI(db *gorm.DB) *RoleDI {
	repo := repository.NewRoleRepository(db)
	usecase := usecase.NewRoleUseCase(repo)
	handler := handler.NewRoleHandler(usecase)

	return &RoleDI{
		db:      db,
		Repo:    repo,
		Usecase: usecase,
		Handler: handler,
	}
}

func (di *RoleDI) Init() *handler.RoleHandler {
	return di.Handler
}
