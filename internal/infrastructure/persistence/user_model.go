package persistence

import (
	"time"

	"user-system/internal/domain/entity"
)

// UserModel GORM 持久化模型
type UserModel struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Username     string     `gorm:"column:username;size:64;uniqueIndex;not null"`
	PasswordHash string     `gorm:"column:password_hash;size:255;not null"`
	Email        string     `gorm:"column:email;size:128"`
	IsLogin      bool       `gorm:"column:is_login;default:false"`
	LoginAt      *time.Time `gorm:"column:login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

// ToEntity 持久化模型 -> 领域实体
func (m *UserModel) ToEntity() *entity.User {
	return &entity.User{
		ID:           m.ID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		Email:        m.Email,
		IsLogin:      m.IsLogin,
		LoginAt:      m.LoginAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// FromEntity 领域实体 -> 持久化模型
func FromEntity(u *entity.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Email:        u.Email,
		IsLogin:      u.IsLogin,
		LoginAt:      u.LoginAt,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
