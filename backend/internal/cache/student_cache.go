package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

type StudentCache struct {
	client *redis.Client
}

func NewStudentCache(client *redis.Client) *StudentCache {
	return &StudentCache{client: client}
}

func (c *StudentCache) SetList(ctx context.Context, key string, students []domain.Student) error {
	newKey := "students"+key

	value, err := json.Marshal(students)
	if err != nil {
		return fmt.Errorf("marshal students: %w", err)
	}

	if err := c.client.Set(ctx, newKey, value, time.Minute).Err(); err != nil {
		return fmt.Errorf("set cached students: %w", err)
	}

	return nil
}

func (c *StudentCache) Set(
	ctx context.Context, key string, student domain.Student,
) error {
	newKey := "student:" + key

	value, err := json.Marshal(student)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}

	if err := c.client.Set(ctx, newKey, value, time.Minute).Err(); err != nil {
		return fmt.Errorf("set cached user: %w", err)
	}

	return nil
}

func (c *StudentCache) GetList(
	ctx context.Context,
	key string,
) ([]domain.Student, bool, error) {
	newKey := "students"+key

	value, err := c.client.Get(ctx, newKey).Result()

	if errors.Is(err, redis.Nil) {
		return []domain.Student{}, false, nil
	}

	if err != nil {
		return []domain.Student{}, false, fmt.Errorf("get cached students: %w", err)
	}

	var result []domain.Student

	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return []domain.Student{}, false, fmt.Errorf("unmarshal cached students: %w", err)
	}

	return result, true, nil
}


func (c *StudentCache) Get(
	ctx context.Context, key string,
) (domain.Student, bool, error) {
	newKey := "student:" + key

	value, err := c.client.Get(ctx, newKey).Result()

	if errors.Is(err, redis.Nil) {
		return domain.Student{}, false, nil
	}

	if err != nil {
		return domain.Student{}, false, fmt.Errorf("get cached user: %w", err)
	}

	var result domain.Student

	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return domain.Student{}, false, fmt.Errorf("unmarshal cached user: %w", err)
	}

	return result, true, nil
}
