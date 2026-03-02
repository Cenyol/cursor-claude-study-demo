package repository

import (
	"context"
	"time"

	"user-system/internal/domain/entity"
)

// UserRepository 由 Domain 定义，Infrastructure 实现（依赖倒置）
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	UpdateLoginState(ctx context.Context, userID int64, isLogin bool, loginAt *time.Time) error
}
