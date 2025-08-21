package di

import (
	"lms-bootcamp/internal/delivery/handler"
	"lms-bootcamp/internal/repository"
	"lms-bootcamp/internal/usecase"

	"gorm.io/gorm"
)

type SessionDI struct {
	db *gorm.DB
	Repo    *repository.SessionRepository
	Usecase *usecase.SessionUseCase
	Handler *handler.SessionHandler
}

func NewSessionDI(db *gorm.DB) *SessionDI {
	repo := repository.NewSessionRepository(db)
	usecase := usecase.NewSessionUseCase(repo)
	handler := handler.NewSessionHandler(*usecase)

	return &SessionDI{
		db:      db,
		Repo:    repo,
		Usecase: usecase,
		Handler: handler,
	}
}

func (di *SessionDI) Init() *handler.SessionHandler {
	return di.Handler
}
