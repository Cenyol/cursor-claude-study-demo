package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户聚合根，充血模型
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Email        string
	IsLogin      bool
	LoginAt      *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

const bcryptCost = 12

// NewUser 工厂方法：注册时创建，密码立即哈希
func NewUser(username, plainPassword, email string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcryptCost)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &User{
		Username:     username,
		PasswordHash: string(hash),
		Email:        email,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// CheckPassword 校验明文密码
func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
	return err == nil
}

// PublicInfo 对外暴露信息（不含密码）
func (u *User) PublicInfo() map[string]interface{} {
	out := map[string]interface{}{
		"id":         u.ID,
		"username":  u.Username,
		"email":     u.Email,
		"created_at": u.CreatedAt,
	}
	if u.LoginAt != nil {
		out["login_at"] = *u.LoginAt
	}
	return out
}
