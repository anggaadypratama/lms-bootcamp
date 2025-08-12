package repository

import (
	"context"
	"lms-bootcamp/internal/domain/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	base *GormRepository[models.UserModel]
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		base: NewGormRepository[models.UserModel](db),
	}
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*models.UserModel, error) {
    // Method 2: Pakai typed filter
    filters := map[string]interface{}{}
    filters["email"] = email
    
    // Atau bisa juga pakai where clause yang lebih explicit
    users, err := r.base.Join("Role").Where("email = ?", email).GetAll(ctx)
    if err != nil {
        return nil, err
    }
    if len(users) == 0 {
        return nil, nil
    }
    return &users[0], nil
}


func (r *AuthRepository) ResetPassword(ctx context.Context, id string, newPassword string) error {
    user, err := r.base.GetById(ctx, id)
    if err != nil {
        return err
    }
    if user == nil {
        return nil
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user.Password = string(hashedPassword)
    return r.base.Update(ctx, user.ID, user)
}