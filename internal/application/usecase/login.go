package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"user-system/internal/application/dto"
	"user-system/internal/domain/repository"
	"user-system/pkg/token"
)

const sessionExpireSeconds = 7200 // 2 小时

// Login 登录用例：校验用户、更新 is_login/login_at、写 Session、返回 Token
func Login(ctx context.Context, req *dto.LoginRequest, userRepo repository.UserRepository, sessionRepo repository.SessionRepository) (*dto.LoginResult, error) {
	user, err := userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidPassword
	}
	now := time.Now()
	if err := userRepo.UpdateLoginState(ctx, user.ID, true, &now); err != nil {
		return nil, err
	}
	tok, err := token.New()
	if err != nil {
		return nil, fmt.Errorf("token new: %w", err)
	}
	userID := strconv.FormatInt(user.ID, 10)
	if err := sessionRepo.Set(ctx, tok, userID, sessionExpireSeconds); err != nil {
		return nil, err
	}
	return &dto.LoginResult{Token: tok}, nil
}
