package usecase

import (
	"context"

	"user-system/internal/application/dto"
	"user-system/internal/domain/entity"
	"user-system/internal/domain/repository"
)

const minPasswordLen = 8

// Register 注册用例：编排流程，业务规则在 Entity
func Register(ctx context.Context, req *dto.RegisterRequest, userRepo repository.UserRepository) (*dto.RegisterResult, error) {
	if len(req.Password) < minPasswordLen {
		return nil, ErrPasswordTooShort
	}
	exists, err := userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}
	user, err := entity.NewUser(req.Username, req.Password, req.Email)
	if err != nil {
		return nil, err
	}
	if err := userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return &dto.RegisterResult{User: user.PublicInfo()}, nil
}
