package usecase

import (
	"context"
	"strconv"

	"user-system/internal/domain/repository"
)

// Logout 退出用例：删 Redis Token，更新用户 is_login/login_at
func Logout(ctx context.Context, tokenStr string, userRepo repository.UserRepository, sessionRepo repository.SessionRepository) error {
	userIDStr, exists, err := sessionRepo.Get(ctx, tokenStr)
	if err != nil {
		return err
	}
	if err := sessionRepo.Delete(ctx, tokenStr); err != nil {
		return err
	}
	if exists && userIDStr != "" {
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)
		_ = userRepo.UpdateLoginState(ctx, userID, false, nil)
	}
	return nil
}
