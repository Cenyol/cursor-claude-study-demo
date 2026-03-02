package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"user-system/internal/domain/repository"
)

const keyPrefix = "user:session:"

type sessionRepository struct {
	client *redis.Client
}

// NewSessionRepository 依赖注入
func NewSessionRepository(client *redis.Client) repository.SessionRepository {
	return &sessionRepository{client: client}
}

func (r *sessionRepository) Set(ctx context.Context, token, userID string, expireSeconds int) error {
	key := keyPrefix + token
	return r.client.Set(ctx, key, userID, time.Duration(expireSeconds)*time.Second).Err()
}

func (r *sessionRepository) Get(ctx context.Context, token string) (userID string, exists bool, err error) {
	key := keyPrefix + token
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

func (r *sessionRepository) Delete(ctx context.Context, token string) error {
	key := keyPrefix + token
	return r.client.Del(ctx, key).Err()
}

var _ repository.SessionRepository = (*sessionRepository)(nil)
