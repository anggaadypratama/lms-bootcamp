package usecase

import (
	"context"
	"fmt"
	"lms-bootcamp/internal/domain/dto"
	"lms-bootcamp/internal/pkg/email"
	"lms-bootcamp/internal/pkg/utils"
	"lms-bootcamp/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repository *repository.AuthRepository
	utils      *utils.Utils
	email 	*email.EmailService
}

func NewAuthUseCase(repo *repository.AuthRepository, utils *utils.Utils) *AuthUseCase {
	emailService := email.NewEmailService()
	return &AuthUseCase{
		repository: repo,
		utils:      utils,
		email:      emailService,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, email string, password string) (*dto.LoginData, error) {
	user, err := uc.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	token, err := uc.utils.GenerateJWTToken(&dto.UserData{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role: &dto.RoleData{
			ID:   user.RoleId,
			Name: user.Role.Name,
		},
	})

	if err != nil {
		return nil, err
	}

	return &dto.LoginData{
		Token: token,
		User: &dto.UserData{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role: &dto.RoleData{
				ID:   user.RoleId,
				Name: user.Role.Name,
			},
		},
	}, nil
}

func (uc *AuthUseCase) ForgotPassword(ctx context.Context, email string) error {
	fmt.Println("Sending forgot password email to:", email)
	user, err := uc.repository.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	token, err := uc.utils.GenerateJWTToken(&dto.UserData{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role: &dto.RoleData{
			ID:   user.RoleId,
			Name: user.Role.Name,
		},
	})

	if err != nil {
		return err
	}

	fmt.Println("Sending forgot password email to:", token)

	err = <-uc.email.SendEmailAsyncWithResult(email, "Forgot Password", token)

	if err != nil {
    fmt.Println("Gagal kirim email:", err)
} else {
    fmt.Println("Email sukses dikirim, kak!")
}

	return err
}

func (uc *AuthUseCase) ResetPassword(ctx context.Context, token string, newPassword string, confirmPassword string) error {
	if newPassword != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	res, err := uc.utils.ValidateJWTToken(token)
	if err != nil {
		return err
	}

	responseData, ok := res["data"].(map[string]interface{})
	if !ok {
		return err
	}

	userId, ok := responseData["id"].(string)
	if !ok {
		return err
	}

	if err := uc.repository.ResetPassword(ctx, userId, newPassword); err != nil {
		return err
	}

	return nil
}