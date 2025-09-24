package usecase

import (
	"lms-bootcamp/internal/pkg/s3"
)

type AssignmentUsecase struct {
	s3 *s3.S3Config
}

func NewAssignmentUsecase(s3 *s3.S3Config) *AssignmentUsecase {
	return &AssignmentUsecase{
		s3: s3,
	}
}
