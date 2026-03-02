package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"user-system/internal/domain/entity"
	"user-system/internal/domain/repository"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 依赖注入
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	m := FromEntity(user)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var m UserModel
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserModel{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) UpdateLoginState(ctx context.Context, userID int64, isLogin bool, loginAt *time.Time) error {
	updates := map[string]interface{}{"is_login": isLogin, "login_at": loginAt}
	return r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", userID).Updates(updates).Error
}
