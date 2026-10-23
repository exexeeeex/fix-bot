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

type GroupCache struct {
	client *redis.Client
}

func NewGroupCache(client *redis.Client) *GroupCache {
	return &GroupCache{client: client}
}

func (c *GroupCache) Set(ctx context.Context, key string, group domain.Group) error {
	newKey := "group" + key

	value, err := json.Marshal(group)
	if err != nil {
		return fmt.Errorf("marshal group: %w", err)
	}

	if err := c.client.Set(ctx, newKey, value, time.Minute).Err(); err != nil {
		return fmt.Errorf("set cached group: %w", err)
	}

	return nil
}

func (c *GroupCache) SetList(ctx context.Context, key string, groups []domain.Group) error {
	newKey := "groups"+key

	value, err := json.Marshal(groups)
	if err != nil {
		return fmt.Errorf("marshal groups: %w", err)
	}

	if err := c.client.Set(ctx, newKey, value, time.Minute).Err(); err != nil {
		return fmt.Errorf("set cached groups: %w", err)
	}

	return nil
}

func (c *GroupCache) Get(
	ctx context.Context, key string,
) (domain.Group, bool, error) {
	newKey := "group" + key

	value, err := c.client.Get(ctx, newKey).Result()

	if errors.Is(err, redis.Nil) {
		return domain.Group{}, false, nil
	}

	if err != nil {
		return domain.Group{}, false, fmt.Errorf("get cached group: %w", err)
	}

	var result domain.Group

	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return domain.Group{}, false, fmt.Errorf("unmarshal group: %w", err)
	}

	return result, true, nil
}

func (c *GroupCache) GetList(
	ctx context.Context,
	key string,
) ([]domain.Group, bool, error) {
	newKey := "groups"+key

	value, err := c.client.Get(ctx, newKey).Result()

	if errors.Is(err, redis.Nil) {
		return []domain.Group{}, false, nil
	}

	if err != nil {
		return []domain.Group{}, false, fmt.Errorf("get cached groups: %w", err)
	}

	var result []domain.Group

	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return []domain.Group{}, false, fmt.Errorf("unmarshal cached groups: %w", err)
	}

	return result, true, nil
}
