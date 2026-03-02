package repository

import "context"

// SessionRepository Token 存储接口，由 Infrastructure 用 Redis 实现
type SessionRepository interface {
	Set(ctx context.Context, token, userID string, expireSeconds int) error
	Get(ctx context.Context, token string) (userID string, exists bool, err error)
	Delete(ctx context.Context, token string) error
}
