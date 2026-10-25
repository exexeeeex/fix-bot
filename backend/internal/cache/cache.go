package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache[T any] struct {
	client *redis.Client
	prefix string
}

func New[T any](
	client *redis.Client,
	prefix string,
) *Cache[T] {
	return &Cache[T]{
		client: client,
		prefix: prefix,
	}
}

func (c *Cache[T]) Set(
	ctx context.Context,
	key string,
	value T,
	ttl int,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("data marshal: %w", err)
	}

	if err := c.client.Set(
		ctx,
		c.prefix+key,
		data,
		time.Minute*time.Duration(ttl),
	).Err(); err != nil {
		return fmt.Errorf("set cache value: %w", err)
	}

	return nil
}

func (c *Cache[T]) Get(
	ctx context.Context,
	key string,
) (T, bool, error) {
	var result T

	data, err := c.client.Get(
		ctx,
		c.prefix+key,
	).Bytes()

	if errors.Is(err, redis.Nil) {
		return result, false, nil
	}

	if err != nil {
		return result, false, fmt.Errorf(
			"get cached value: %w", err,
		)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, false, fmt.Errorf(
			"unmarshal cache value: %w", err,
		)
	}

	return result, true, nil
}

func (c *Cache[T]) Delete(
	ctx context.Context,
	key string,
) error {
	if err := c.client.Del(
		ctx,
		c.prefix+key,
	).Err(); err != nil {
		return fmt.Errorf("delete cache value: %w", err)
	}

	return nil
}
