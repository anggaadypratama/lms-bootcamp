package handler

import (
	"lms-bootcamp/internal/usecase"
)

type AssignmentHandler struct {
	usecase *usecase.AssignmentUsecase
}

func NewAssignmentHandler(usecase *usecase.AssignmentUsecase) *AssignmentHandler {
	return &AssignmentHandler{
		usecase: usecase,
	}
}
