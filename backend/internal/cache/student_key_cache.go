package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StudentKeyCache struct {
	client *redis.Client
}

func NewStudentKeyCache(client *redis.Client) *StudentKeyCache {
	return &StudentKeyCache{client: client}
}

func (c *StudentKeyCache) Set(
	ctx context.Context,
	key string,
	studentKey string,
) error {
	newKey := "student-key:" + key

	if err := c.client.Set(ctx, newKey, studentKey, time.Minute*10).Err(); err != nil {
		return fmt.Errorf("set cached key: %w", err)
	}

	return nil
}

func (c *StudentKeyCache) Get(
	ctx context.Context,
	key string,
) (string, bool, error) {
	newKey := "student-key:" + key

	result, err := c.client.Get(ctx, newKey).Result()

	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("get cached key: %w", err)
	}

	return result, true, nil
}
