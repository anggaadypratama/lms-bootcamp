package usecase

import (
	"context"
)

type AuthUseCase interface {
	Login(ctx context.Context, email string, password string) (string, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
}