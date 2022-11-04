package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Manager struct {
	client *redis.Client
}

func NewManager(client *redis.Client) *Manager {
	return &Manager{
		client: client,
	}
}

func (m *Manager) Set(ctx context.Context, key string, value interface{}, tags []string, expiration time.Duration) error {
	pipe := m.client.TxPipeline()
	for _, tag := range tags {
		pipe.SAdd(ctx, tag, key)
		pipe.Expire(ctx, tag, expiration)
	}

	pipe.Set(ctx, key, value, expiration)

	_, errExec := pipe.Exec(ctx)
	return errExec
}

func (m *Manager) Invalidate(ctx context.Context, tags []string) {
	keys := make([]string, 0)
	for _, tag := range tags {
		k, _ := m.client.SMembers(ctx, tag).Result()
		keys = append(keys, tag)
		keys = append(keys, k...)
	}
	m.client.Del(ctx, keys...)
}

func (m *Manager) GetBytes(ctx context.Context, key string) ([]byte, error) {
	data, err := m.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, errNotFoundInCache
		}
		return nil, errors.Wrap(err, "cannot get data from cache")
	}
	return data, nil
}
